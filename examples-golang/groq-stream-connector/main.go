// Groq Stream Connector - Streaming Chat Completions
// Uses Vastar Connector SDK (Simulator + Real Groq API Mode)
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
)

//
// =======================
// OpenAI-compatible Types (Groq compatible)
// =======================
//

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type StreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

//
// =======================
// Connector Definition
// =======================
//

type GroqConnector struct {
	client  *vastar.RuntimeClient
	baseURL string
	apiKey  string
}

func NewGroqConnector(baseURL, apiKey string) (*GroqConnector, error) {
	client, err := vastar.NewRuntimeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create runtime client: %w", err)
	}

	return &GroqConnector{
		client:  client,
		baseURL: baseURL,
		apiKey:  apiKey,
	}, nil
}

func (c *GroqConnector) Close() error {
	return c.client.Close()
}

//
// =======================
// Test Connection (Simulator Only)
// =======================
//

func (c *GroqConnector) TestConnection() error {
	req := vastar.POST(c.baseURL + "/test_completion").
		WithHeader("Content-Type", "application/json").
		WithTimeout(30_000)

	resp, err := c.client.ExecuteHTTP(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

//
// =======================
// Streaming Chat Completion
// =======================
//

func (c *GroqConnector) ChatCompletionStream(
	req ChatCompletionRequest,
) (<-chan string, <-chan error) {

	out := make(chan string, 100)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)

		req.Stream = true

		body, err := json.Marshal(req)
		if err != nil {
			errCh <- err
			return
		}

		httpReq := vastar.POST(c.baseURL + "/v1/chat/completions").
			WithHeader("Content-Type", "application/json").
			WithHeader("Accept", "text/event-stream").
			WithBody(body).
			WithTimeout(300_000)

		if c.apiKey != "" {
			httpReq = httpReq.WithHeader("Authorization", "Bearer "+c.apiKey)
		}

		resp, err := c.client.ExecuteHTTP(httpReq)
		if err != nil {
			errCh <- err
			return
		}

		if resp.StatusCode != 200 {
			errCh <- fmt.Errorf(
				"unexpected status code: %d - %s",
				resp.StatusCode,
				string(resp.Body),
			)
			return
		}

		scanner := bufio.NewScanner(bytes.NewReader(resp.Body))
		for scanner.Scan() {
			line := scanner.Text()

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}

			var chunk StreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			for _, choice := range chunk.Choices {
				if choice.Delta.Content != "" {
					out <- choice.Delta.Content
				}
				if choice.FinishReason != nil {
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			errCh <- err
		}
	}()

	return out, errCh
}

//
// =======================
// Non-streaming Helper
// =======================
//

func (c *GroqConnector) ChatCompletion(req ChatCompletionRequest) (string, error) {
	chunks, errs := c.ChatCompletionStream(req)
	var full strings.Builder

	for {
		select {
		case chunk, ok := <-chunks:
			if !ok {
				return full.String(), nil
			}
			full.WriteString(chunk)

		case err := <-errs:
			if err != nil {
				return "", err
			}
		}
	}
}

//
// =======================
// Main (Mentor Style)
// =======================
//

func main() {
	fmt.Println("🤖 Groq Stream Connector")
	fmt.Println(strings.Repeat("=", 60))

	apiKey := os.Getenv("GROQ_API_KEY")
	baseURL := os.Getenv("GROQ_BASE_URL")

	// ===== Mode Selection =====
	if baseURL == "" {
		if apiKey != "" {
			baseURL = "https://api.groq.com/openai"
			fmt.Println("Mode  : REAL GROQ API")
		} else {
			baseURL = "http://localhost:4545"
			fmt.Println("Mode  : RAI SIMULATOR")
		}
	}

	fmt.Println("Base  :", baseURL)
	fmt.Println(strings.Repeat("-", 60))

	connector, err := NewGroqConnector(baseURL, apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer connector.Close()

	// ===== Simulator test =====
	if apiKey == "" {
		fmt.Println("Testing simulator connection...")
		if err := connector.TestConnection(); err != nil {
			log.Fatal("Simulator not reachable:", err)
		}
		fmt.Println("Simulator OK ✅")
		fmt.Println(strings.Repeat("-", 60))
	}

	// ===== Example: Streaming =====
	req := ChatCompletionRequest{
		Model: "llama-3.1-8b-instant", // bebas di simulator
		Messages: []Message{
			{
				Role:    "user",
				Content: "Explain quantum computing in simple terms.",
			},
		},
		Temperature: 0.7,
		MaxTokens:   300,
	}

	fmt.Println("AI Response:")
	fmt.Println(strings.Repeat("-", 60))

	start := time.Now()

	chunks, errs := connector.ChatCompletionStream(req)
	var chars int

	for {
		select {
		case c, ok := <-chunks:
			if !ok {
				fmt.Println()
				fmt.Println(strings.Repeat("-", 60))
				fmt.Printf("Completed in %v\n", time.Since(start))
				fmt.Printf("Characters : %d\n", chars)
				fmt.Println(strings.Repeat("=", 60))
				return
			}
			fmt.Print(c)
			chars += len(c)

		case err := <-errs:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
