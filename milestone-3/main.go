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

	// Read-Only promp
	// prompt := "Read the file demo.txt, and count the number of characters"

	// Read-Write prompt
	prompt := "Read the file demo.txt, then write its contents reversed to reversed.txt"

	messages := []Message{{Role: "user", Content: prompt}}

	for {
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

		if len(msg.ToolCalls) == 0 {
			fmt.Println(msg.Content)
			break
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
}
