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
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
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

	prompt := "What is 2+2?"

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

	fmt.Println(result.Message.Content)
}
