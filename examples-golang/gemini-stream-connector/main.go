// Gemini Stream Connector - Streaming Generate Content
// Uses Vastar Connector SDK (Simulator + Real API Mode)
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
// Gemini API Types
// =======================
//

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Role  string `json:"role"` // user | model | system
	Parts []Part `json:"parts"`
}

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type GeminiStreamChunk struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

//
// =======================
// OpenAI Simulator Types
// =======================
//

type OpenAIStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

//
// =======================
// Connector Definition
// =======================
//

type GeminiConnector struct {
	client  *vastar.RuntimeClient
	baseURL string
	apiKey  string
	model   string
}

func NewGeminiConnector(baseURL, apiKey, model string) (*GeminiConnector, error) {
	client, err := vastar.NewRuntimeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create runtime client: %w", err)
	}

	return &GeminiConnector{
		client:  client,
		baseURL: baseURL,
		apiKey:  apiKey,
		model:   model,
	}, nil
}

func (c *GeminiConnector) Close() error {
	return c.client.Close()
}

//
// =======================
// Test Connection
// =======================
//

func (c *GeminiConnector) TestConnection() error {
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
// Streaming GenerateContent
// =======================
//

func (c *GeminiConnector) StreamGenerateContent(req GeminiRequest) (<-chan string, <-chan error) {
	out := make(chan string, 100)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)

		// ======================
		// REAL GEMINI MODE
		// ======================
		if c.apiKey != "" {
			body, _ := json.Marshal(req)

			url := fmt.Sprintf(
				"%s/v1beta/models/%s:streamGenerateContent?key=%s",
				c.baseURL,
				c.model,
				c.apiKey,
			)

			httpReq := vastar.POST(url).
				WithHeader("Content-Type", "application/json").
				WithHeader("Accept", "text/event-stream").
				WithBody(body).
				WithTimeout(300_000)

			resp, err := c.client.ExecuteHTTP(httpReq)
			if err != nil {
				errCh <- err
				return
			}

			if resp.StatusCode != 200 {
				errCh <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				return
			}

			scanner := bufio.NewScanner(bytes.NewReader(resp.Body))
			for scanner.Scan() {
				line := scanner.Text()
				if !strings.HasPrefix(line, "data: ") {
					continue
				}

				data := strings.TrimPrefix(line, "data: ")
				var chunk GeminiStreamChunk
				if json.Unmarshal([]byte(data), &chunk) != nil {
					continue
				}

				for _, cand := range chunk.Candidates {
					for _, part := range cand.Content.Parts {
						if part.Text != "" {
							out <- part.Text
						}
					}
				}
			}
			return
		}

		// ======================
		// SIMULATOR MODE
		// ======================
		openAIReq := map[string]interface{}{
			"model": "gpt-4",
			"messages": []map[string]string{
				{"role": "user", "content": req.Contents[0].Parts[0].Text},
			},
			"stream": true,
		}

		body, _ := json.Marshal(openAIReq)

		httpReq := vastar.POST(c.baseURL + "/v1/chat/completions").
			WithHeader("Content-Type", "application/json").
			WithHeader("Accept", "text/event-stream").
			WithBody(body).
			WithTimeout(300_000)

		resp, err := c.client.ExecuteHTTP(httpReq)
		if err != nil {
			errCh <- err
			return
		}

		if resp.StatusCode != 200 {
			errCh <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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

			var chunk OpenAIStreamChunk
			if json.Unmarshal([]byte(data), &chunk) != nil {
				continue
			}

			for _, c := range chunk.Choices {
				if c.Delta.Content != "" {
					out <- c.Delta.Content
				}
			}
		}
	}()

	return out, errCh
}

//
// =======================
// Main
// =======================
//

func main() {
	fmt.Println("🤖 Gemini Stream Connector")
	fmt.Println(strings.Repeat("=", 60))

	apiKey := os.Getenv("GEMINI_API_KEY")
	baseURL := os.Getenv("GEMINI_BASE_URL")

	if baseURL == "" {
		if apiKey != "" {
			baseURL = "https://generativelanguage.googleapis.com"
			fmt.Println("Mode  : REAL GEMINI API")
		} else {
			baseURL = "http://localhost:8080"
			fmt.Println("Mode  : OPENAI SIMULATOR")
		}
	}

	fmt.Println("Base  :", baseURL)
	fmt.Println(strings.Repeat("-", 60))

	connector, err := NewGeminiConnector(baseURL, apiKey, "gemini-2.0-flash")
	if err != nil {
		log.Fatal(err)
	}
	defer connector.Close()

	if apiKey == "" {
		if err := connector.TestConnection(); err != nil {
			log.Fatal("Simulator not reachable:", err)
		}
		fmt.Println("Status: Simulator connected ✅")
		fmt.Println(strings.Repeat("-", 60))
	}

	req := GeminiRequest{
		Contents: []Content{
			{
				Role: "user",
				Parts: []Part{
					{Text: "Explain quantum computing in simple terms."},
				},
			},
		},
	}

	fmt.Println("AI Response:")
	fmt.Println(strings.Repeat("-", 60))

	start := time.Now()

	chunks, errors := connector.StreamGenerateContent(req)
	var totalChars int

	for {
		select {
		case c, ok := <-chunks:
			if !ok {
				duration := time.Since(start)
				fmt.Println()
				fmt.Println(strings.Repeat("-", 60))
				fmt.Printf("Completed in %v\n", duration)
				fmt.Printf("Characters : %d\n", totalChars)
				fmt.Println(strings.Repeat("=", 60))
				return
			}
			fmt.Print(c)
			totalChars += len(c)

		case err := <-errors:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
