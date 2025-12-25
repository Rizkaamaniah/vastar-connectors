/**
 * VASTAR GROQ CONNECTOR - PRODUCTION SAFE EDITION
 * Version: 1.1.0-Stable
 *
 * NOTE:
 * - Buffered streaming (NOT real-time SSE)
 * - Fully compatible with Vastar Runtime IPC
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
	ipc "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang/protocol"
)

////////////////////////////////////////////////////////////////////////////////
// CONFIGURATION
////////////////////////////////////////////////////////////////////////////////

const (
	GroqBaseURL      = "https://api.groq.com/openai/v1/chat/completions"
	DefaultGroqModel = "llama-3.1-70b-versatile"

	RequestTimeoutMS = 300000 // 5 menit (enterprise-safe)
)

////////////////////////////////////////////////////////////////////////////////
// GROQ API CONTRACT
////////////////////////////////////////////////////////////////////////////////

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Model       string        `json:"model"`
	Messages    []GroqMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
	Stream      bool          `json:"stream"`
}

/**
 * Groq mengembalikan ARRAY of chunks (buffered by runtime)
 */
type GroqStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason,omitempty"`
	} `json:"choices"`
}

////////////////////////////////////////////////////////////////////////////////
// CONNECTOR CORE
////////////////////////////////////////////////////////////////////////////////

type GroqConnector struct {
	client *vastar.RuntimeClient
	apiKey string
	logger *log.Logger
}

func NewGroqConnector(apiKey string) (*GroqConnector, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GROQ_API_KEY is not set")
	}

	client, err := vastar.NewRuntimeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect runtime: %w", err)
	}

	return &GroqConnector{
		client: client,
		apiKey: apiKey,
		logger: log.New(os.Stdout, "[GROQ] ", log.LstdFlags),
	}, nil
}

func (c *GroqConnector) Close() error {
	c.logger.Println("closing IPC connection")
	return c.client.Close()
}

////////////////////////////////////////////////////////////////////////////////
// EXECUTION (BUFFERED STREAM)
////////////////////////////////////////////////////////////////////////////////

func (c *GroqConnector) Execute(prompt string) error {
	reqBody := GroqRequest{
		Model: DefaultGroqModel,
		Messages: []GroqMessage{
			{Role: "user", Content: prompt},
		},
		Temperature: 0.7,
		Stream:      true,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req := vastar.POST(GroqBaseURL).
		WithHeader("Authorization", "Bearer "+c.apiKey).
		WithHeader("Content-Type", "application/json").
		WithBody(payload).
		WithTimeout(RequestTimeoutMS)

	resp, err := c.client.ExecuteHTTP(req)
	if err != nil {
		return err
	}

	if resp.ErrorClass != ipc.ErrorClassSuccess {
		return fmt.Errorf("runtime error class: %s", resp.ErrorClass)
	}

	var chunks []GroqStreamChunk
	if err := json.Unmarshal(resp.Body, &chunks); err != nil {
		return fmt.Errorf("invalid groq stream payload")
	}

	for _, chunk := range chunks {
		if len(chunk.Choices) == 0 {
			continue
		}

		content := chunk.Choices[0].Delta.Content
		if content != "" {
			fmt.Print(content)
		}

		if chunk.Choices[0].FinishReason != nil {
			break
		}
	}

	fmt.Println()
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN ENTRYPOINT
////////////////////////////////////////////////////////////////////////////////

func main() {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		log.Fatal("GROQ_API_KEY environment variable missing")
	}

	connector, err := NewGroqConnector(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer connector.Close()

	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("VASTAR GROQ CONNECTOR (PRODUCTION SAFE)")
	fmt.Println("Execution : Buffered Stream via IPC")
	fmt.Println("Protocol  : FlatBuffers / Runtime Managed")
	fmt.Println(strings.Repeat("=", 70))

	prompt := "Explain why low-latency IPC is important in modern workflow engines."

	fmt.Print("🤖 Groq: ")
	start := time.Now()

	if err := connector.Execute(prompt); err != nil {
		log.Fatal("execution failed:", err)
	}

	fmt.Println("Latency:", time.Since(start))
	fmt.Println(strings.Repeat("=", 70))
}
