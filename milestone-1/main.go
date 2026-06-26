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
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ChatResponse struct {
	Message Message `json:"message"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	prompt := "What is 2+2?"

	chatRequest := ChatRequest{
		Model: "llama3.2",
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		Stream: false,
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
