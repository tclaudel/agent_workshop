package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
	Tools    []ToolDef `json:"tools,omitempty"`
}

type ChatResponse struct {
	Message Message `json:"message"`
}

type ToolCallFunction struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type ToolCall struct {
	Function ToolCallFunction `json:"function"`
}

type ToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

type ToolDef struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type Tool struct {
	Name        string
	Description string
	Parameters  json.RawMessage
	Run         func(args map[string]any) string
}

var httpClient = &http.Client{Timeout: 60 * time.Second}

// maxTurns caps the agent loop so a model that never converges can't run
// forever. The validation step below can also send the agent back around, so
// this bound protects against an infinite "fix it / still wrong" cycle.
const maxTurns = 10

// checkUppercase is the test we run on the model's output. An LLM is
// probabilistic — "uppercase this text" is a request, not a guarantee — so we
// assert the property we care about: the output must equal the original with
// every letter uppercased. Unlike a general transform, uppercasing is cheap to
// verify *exactly*, so we do. An empty slice means the output passed. The
// slice return lets you add more assertions later and report them all at once.
func checkUppercase(original, output string) []string {
	var failures []string
	if want := strings.ToUpper(original); output != want {
		failures = append(failures, fmt.Sprintf("not fully uppercased: expected %q, got %q", want, output))
	}
	return failures
}

// validateUppercase reads the source and result files off disk and runs
// checkUppercase on them. It returns ok=false plus a human-readable message
// that we feed straight back to the model so it can correct itself.
func validateUppercase(srcPath, dstPath string) (ok bool, feedback string) {
	src, err := os.ReadFile(srcPath)
	if err != nil {
		return false, fmt.Sprintf("could not read source %s: %v", srcPath, err)
	}
	dst, err := os.ReadFile(dstPath)
	if err != nil {
		return false, fmt.Sprintf("could not read %s — did you write it? %v", dstPath, err)
	}

	// Ignore trailing line endings (\n or \r\n) on either side so the check
	// focuses on the content the model is responsible for, not file-ending
	// conventions — demo.txt may have Windows CRLF.
	original := strings.TrimRight(string(src), "\r\n")
	output := strings.TrimRight(string(dst), "\r\n")

	if failures := checkUppercase(original, output); len(failures) > 0 {
		return false, "validation failed: " + strings.Join(failures, "; ")
	}
	return true, ""
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	tools := []Tool{
		{
			Name:        "read_file",
			Description: "Read the contents of a file",
			Parameters:  json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"}},"required":["path"]}`),
			Run: func(args map[string]any) string {
				path, ok := args["path"].(string)
				if !ok || path == "" {
					return "error: missing or invalid 'path' argument"
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return fmt.Sprintf("error: %v", err)
				}
				return string(data)
			},
		},
		{
			Name:        "write_file",
			Description: "Write content to a file",
			Parameters:  json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"},"content":{"type":"string"}},"required":["path","content"]}`),
			Run: func(args map[string]any) string {
				path, ok := args["path"].(string)
				if !ok || path == "" {
					return "error: missing or invalid 'path' argument"
				}
				content, _ := args["content"].(string)
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					return fmt.Sprintf("error: %v", err)
				}
				return "ok"
			},
		},
	}

	toolsByName := make(map[string]*Tool, len(tools))
	for i := range tools {
		toolsByName[tools[i].Name] = &tools[i]
	}

	toolDefs := make([]ToolDef, len(tools))
	for i, t := range tools {
		toolDefs[i] = ToolDef{
			Type: "function",
			Function: ToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  t.Parameters,
			},
		}
	}

	prompt := "Read the file demo.txt, convert its contents to UPPERCASE, then write the result to uppercase.txt"

	messages := []Message{{Role: "user", Content: prompt}}

	for turn := 0; turn < maxTurns; turn++ {
		chatRequest := ChatRequest{
			Model:    "llama3.2",
			Messages: messages,
			Stream:   false,
			Tools:    toolDefs,
		}

		body, err := json.Marshal(chatRequest)
		if err != nil {
			log.Fatal(err)
		}

		logger.Info("Send Message", slog.Any("body", chatRequest))

		resp, err := httpClient.Post("http://localhost:11434/api/chat", "application/json", bytes.NewReader(body))
		if err != nil {
			log.Fatal(err)
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("ollama error %d: %s", resp.StatusCode, data)
		}

		var result ChatResponse
		if err := json.Unmarshal(data, &result); err != nil {
			log.Fatal(err)
		}

		logger.Info("Received Response", slog.Any("Response", result))

		msg := result.Message
		messages = append(messages, msg)

		logger.Info("Tools", slog.Any("tool_calls", msg.ToolCalls))

		if len(msg.ToolCalls) == 0 {
			// The model thinks it is done. Don't trust it — verify the output
			// against our deterministic checks first.
			ok, feedback := validateUppercase("demo.txt", "uppercase.txt")
			if ok {
				logger.Info("validation passed")
				fmt.Println(msg.Content)
				return
			}
			// Tests failed. Hand the failure back as a new user turn and let
			// the agent loop take another shot. This is the loop self-correcting.
			logger.Warn("validation failed", slog.String("feedback", feedback))
			messages = append(messages, Message{Role: "user", Content: feedback + " — please fix uppercase.txt."})
			continue
		}

		for _, tc := range msg.ToolCalls {
			tool, ok := toolsByName[tc.Function.Name]
			if !ok {
				messages = append(messages, Message{Role: "tool", Content: fmt.Sprintf("error: unknown tool %q", tc.Function.Name)})
				continue
			}
			output := tool.Run(tc.Function.Arguments)
			messages = append(messages, Message{Role: "tool", Content: output})
		}
	}

	log.Fatalf("gave up after %d turns without passing validation", maxTurns)
}
