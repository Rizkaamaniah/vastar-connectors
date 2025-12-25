package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	

	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
	ipc "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang/protocol"
)

////////////////////////////////////////////////////////////////////////////////
// CONFIGURATION
////////////////////////////////////////////////////////////////////////////////

const (
	ConnectorName = "gemini-stream-connector"
	BaseURL       = "https://generativelanguage.googleapis.com/v1beta"
	DefaultModel  = "gemini-2.0-flash"

	RequestTimeoutMS = 120000 // 2 menit (aman untuk workflow)
)

////////////////////////////////////////////////////////////////////////////////
// GEMINI API CONTRACT
////////////////////////////////////////////////////////////////////////////////

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiStreamChunk struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
		FinishReason string `json:"finishReason,omitempty"`
	} `json:"candidates"`
}

////////////////////////////////////////////////////////////////////////////////
// CONNECTOR CORE
////////////////////////////////////////////////////////////////////////////////

type GeminiConnector struct {
	client *vastar.RuntimeClient
	apiKey string
	model  string
	logger *log.Logger
}

func NewGeminiConnector(apiKey string) (*GeminiConnector, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}

	client, err := vastar.NewRuntimeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect runtime: %w", err)
	}

	return &GeminiConnector{
		client: client,
		apiKey: apiKey,
		model:  DefaultModel,
		logger: log.New(os.Stdout, "[GEMINI] ", log.LstdFlags),
	}, nil
}

func (c *GeminiConnector) Close() error {
	c.logger.Println("closing IPC connection")
	return c.client.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STREAM EXECUTION (PRODUCTION SAFE)
////////////////////////////////////////////////////////////////////////////////

func (c *GeminiConnector) Execute(prompt string) error {
	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Role: "user",
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/models/%s:streamGenerateContent?key=%s",
		BaseURL,
		c.model,
		c.apiKey,
	)

	req := vastar.POST(url).
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

	// Gemini mengembalikan ARRAY of stream chunks
	var chunks []GeminiStreamChunk
	if err := json.Unmarshal(resp.Body, &chunks); err != nil {
		return fmt.Errorf("invalid gemini stream payload")
	}

	for _, chunk := range chunks {
		if len(chunk.Candidates) == 0 {
			continue
		}

		candidate := chunk.Candidates[0]

		if len(candidate.Content.Parts) > 0 {
			fmt.Print(candidate.Content.Parts[0].Text)
		}

		if candidate.FinishReason != "" {
			break
		}
	}

	fmt.Println()
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// MAIN ENTRYPOINT (WORKFLOW READY)
////////////////////////////////////////////////////////////////////////////////

func main() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable missing")
	}

	connector, err := NewGeminiConnector(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer connector.Close()

	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("VASTAR GEMINI STREAM CONNECTOR (PRODUCTION)")
	fmt.Println("Protocol : IPC / FlatBuffers")
	fmt.Println("Runtime  : Managed")
	fmt.Println(strings.Repeat("=", 70))

	prompt := "Explain why enterprises prefer IPC over HTTP for workflow orchestration."

	fmt.Print("🤖 Gemini: ")
	if err := connector.Execute(prompt); err != nil {
		log.Fatal("execution failed:", err)
	}

	fmt.Println("\nStatus:", ipc.ErrorClassSuccess.String())
	fmt.Println(strings.Repeat("=", 70))
}
