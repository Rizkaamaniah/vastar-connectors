/**
 * VASTAR CLAUDE STREAM CONNECTOR - ENTERPRISE EDITION
 * Version: 2.1.0-Stable
 * -----------------------------------------------------------------------------
 * Overview:
 * Konektor ini mengintegrasikan Anthropic Claude 3.5 Sonnet ke dalam 
 * ekosistem Vastar menggunakan protokol biner FlatBuffers.
 *
 * Arsitektur:
 * [App] <--- IPC (Unix Socket) ---> [Vastar Runtime] <--- SSE (HTTPS) ---> [Claude API]
 * -----------------------------------------------------------------------------
 */

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

	// Vastar SDK Core & IPC Protocols
	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
)

// =============================================================================
// 1. DATA MODELS & ANTHROPIC PROTOCOLS
// =============================================================================

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []AnthropicMessage `json:"messages"`
	Stream    bool               `json:"stream"`
}

type AnthropicStreamEvent struct {
	Type  string `json:"type"`
	Delta struct {
		Text string `json:"text"`
	} `json:"delta"`
}

// =============================================================================
// 2. CLAUDE CONNECTOR ENGINE
// =============================================================================

type ClaudeConnector struct {
	client  *vastar.RuntimeClient
	apiKey  string
	version string
}

// NewClaudeConnector menginisialisasi client dengan handshake Unix Socket
func NewClaudeConnector(apiKey string) (*ClaudeConnector, error) {
	client, err := vastar.NewRuntimeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to reach vastar runtime: %w", err)
	}
	return &ClaudeConnector{
		client:  client,
		apiKey:  apiKey,
		version: "2023-06-01",
	}, nil
}

// ChatStream melakukan pemrosesan streaming biner melalui Vastar IPC
func (c *ClaudeConnector) ChatStream(prompt string) (<-chan string, <-chan error) {
	chunkChan := make(chan string, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(chunkChan)
		defer close(errChan)

		reqBody := AnthropicRequest{
			Model:     "claude-3-5-sonnet-20240620",
			MaxTokens: 1024,
			Messages:  []AnthropicMessage{{Role: "user", Content: prompt}},
			Stream:    true,
		}
		body, _ := json.Marshal(reqBody)

		// Routing request melalui Vastar Runtime via FlatBuffers IPC
		httpReq := vastar.POST("https://api.anthropic.com/v1/messages").
			WithHeader("x-api-key", c.apiKey).
			WithHeader("anthropic-version", c.version).
			WithHeader("content-type", "application/json").
			WithBody(body)

		resp, err := c.client.ExecuteHTTP(httpReq)
		if err != nil {
			errChan <- err
			return
		}

		if resp.StatusCode != 200 {
			errChan <- fmt.Errorf("claude API Error %d: %s", resp.StatusCode, string(resp.Body))
			return
		}

		// Parsing Server-Sent Events (SSE) Stream
		scanner := bufio.NewScanner(bytes.NewReader(resp.Body))
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			
			data := strings.TrimPrefix(line, "data: ")
			var event AnthropicStreamEvent
			if err := json.Unmarshal([]byte(data), &event); err == nil {
				if event.Type == "content_block_delta" {
					chunkChan <- event.Delta.Text
				}
			}
		}
	}()
	return chunkChan, errChan
}

func (c *ClaudeConnector) Close() error {
	return c.client.Close()
}

// =============================================================================
// 3. MAIN INTERACTIVE EXECUTION
// =============================================================================

func main() {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("❌ Error: ANTHROPIC_API_KEY environment variable is not set")
	}

	// TAMPILAN OUTPUT BERSIH (ENTERPRISE STYLE)
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("%54s\n", "VASTAR SDK - CLAUDE HIGH-PERFORMANCE CONNECTOR")
	fmt.Println(strings.Repeat("=", 80))

	connector, err := NewClaudeConnector(apiKey)
	if err != nil {
		log.Fatalf("Initialization Failed: %v", err)
	}
	defer connector.Close()

	// Identitas Koneksi IPC
	fmt.Println("🐧 Connected via Unix Socket: /tmp/vastar-connector-runtime.sock")
	
	prompt := "Explain the importance of low-latency IPC in modern microservices."
	fmt.Printf("User: %s\n", prompt)
	fmt.Printf("Claude AI (Sonnet): ")

	start := time.Now()
	charCount := 0

	// Mengonsumsi stream dari channel
	cChan, eChan := connector.ChatStream(prompt)
	for {
		select {
		case chunk, ok := <-cChan:
			if !ok {
				goto Finalize
			}
			fmt.Print(chunk)
			charCount += len(chunk)
			os.Stdout.Sync()
		case err := <-eChan:
			if err != nil {
				fmt.Printf("\n\n❌ Stream Error: %v\n", err)
				return
			}
		}
	}

Finalize:
	elapsed := time.Since(start)
	fmt.Println("\n\n" + strings.Repeat("-", 80))
	// METRIK AUDIT PERFORMA
	fmt.Printf("📊 Metrics: %d characters generated in %v\n", charCount, elapsed)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("[%s] %s Releasing Claude connection...\n", "CLAUDE-SONNET", time.Now().Format("2006/01/02 15:04:05"))
}