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
	Arguments map[string]any `json:"arguments"` // Ollama returns parsed args as a JSON object
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

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	tools := []Tool{
		{
			Name:        "read_file",
			Description: "Read the contents of a file",
			Parameters:  json.RawMessage(`{"type":"object","properties":{"path":{"type":"string"}},"required":["path"]}`),
			Run: func(args map[string]any) string {
				path, _ := args["path"].(string)
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
				path, _ := args["path"].(string)
				content, _ := args["content"].(string)
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					return fmt.Sprintf("error: %v", err)
				}
				return "ok"
			},
		},
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

	prompt := "Read the file demo.txt and tell me what it contains."

	chatRequest := ChatRequest{
		Model: "llama3.2",
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		Stream: false,
		Tools:  toolDefs,
	}

	body, err := json.Marshal(chatRequest)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("sending request", slog.Any("request", chatRequest))

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result ChatResponse
	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatal(err)
	}

	logger.Info("received response", slog.Any("response", result))

	// The model answered our tool-aware request. Two outcomes:
	//   - plain text  → Content is non-empty, we print it (just like Milestone 1).
	//   - tool call    → the model asks us to run read_file. Content is usually
	//                     empty; the real payload is in ToolCalls.
	// This milestone only DECLARES tools, so we surface the request and stop.
	// Executing it (and feeding the result back) is Milestone 3.
	if len(result.Message.ToolCalls) > 0 {
		for _, tc := range result.Message.ToolCalls {
			fmt.Printf("model wants to call %s(%v) — but Milestone 2 does not run it\n",
				tc.Function.Name, tc.Function.Arguments)
		}
		return
	}

	fmt.Println(result.Message.Content)
}
