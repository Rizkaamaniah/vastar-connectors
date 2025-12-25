/**
 * VASTAR MOCK LLM CONNECTOR - PRODUCTION FINAL
 * Version: 2.5.0-Production
 *
 * Purpose:
 * - CI / Unit Test
 * - Workflow Development
 * - Offline Simulation
 *
 * NO external API
 * NO Runtime IPC dependency
 */

package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	
)

////////////////////////////////////////////////////////////////////////////////
// CONFIGURATION
////////////////////////////////////////////////////////////////////////////////

const (
	DefaultModel = "vastar-mock-v2.5"

	MinTTFTMs      = 300  // Time To First Token (ms)
	MaxTTFTMs      = 800
	MinTokenDelay  = 40   // per-token delay (ms)
	MaxTokenDelay  = 120
	ErrorRatePct   = 5    // 5% simulated failure
	MaxTokenFactor = 2    // token estimation multiplier
)

////////////////////////////////////////////////////////////////////////////////
// DATA MODELS (OPENAI-LIKE)
////////////////////////////////////////////////////////////////////////////////

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

type UsageStats struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE (PRODUCTION SAFE)
////////////////////////////////////////////////////////////////////////////////

type MockLLM interface {
	ChatCompletion(req ChatCompletionRequest) (string, UsageStats, error)
	ChatCompletionStream(req ChatCompletionRequest) (<-chan string, <-chan error)
	Close() error
}

////////////////////////////////////////////////////////////////////////////////
// CONNECTOR IMPLEMENTATION
////////////////////////////////////////////////////////////////////////////////

type MockLLMConnector struct {
	model  string
	logger *log.Logger
}

func NewMockLLMConnector(model string) *MockLLMConnector {
	if model == "" {
		model = DefaultModel
	}

	return &MockLLMConnector{
		model:  model,
		logger: log.New(os.Stdout, "\033[32m[MOCK-LLM]\033[0m ", log.LstdFlags),
	}
}

func (c *MockLLMConnector) Close() error {
	c.logger.Println("Mock engine shutdown completed")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INTERNAL HELPERS
////////////////////////////////////////////////////////////////////////////////

func simulateFailure() error {
	if rand.Intn(100) < ErrorRatePct {
		return errors.New("simulated rate limit / transient error")
	}
	return nil
}

func estimateTokens(text string) int {
	return len(strings.Fields(text)) * MaxTokenFactor
}

func generateMockResponse(prompt string) string {
	p := strings.ToLower(prompt)

	switch {
	case strings.Contains(p, "vastar"):
		return "Vastar Workflow SDK adalah sistem IPC berbasis FlatBuffers yang dirancang untuk performa tinggi dan workflow enterprise."
	case strings.Contains(p, "golang"):
		return "Go adalah bahasa pemrograman modern dengan garbage collector efisien dan concurrency model berbasis goroutine."
	case strings.Contains(p, "ipc"):
		return "IPC memungkinkan komunikasi antar proses dengan latensi rendah dibanding HTTP tradisional."
	default:
		return "Sebagai Mock AI, saya membantu Anda menguji workflow tanpa ketergantungan API eksternal."
	}
}

////////////////////////////////////////////////////////////////////////////////
// STREAMING MODE
////////////////////////////////////////////////////////////////////////////////

func (c *MockLLMConnector) ChatCompletionStream(
	req ChatCompletionRequest,
) (<-chan string, <-chan error) {

	textChan := make(chan string, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(textChan)
		defer close(errChan)

		if err := simulateFailure(); err != nil {
			errChan <- err
			return
		}

		time.Sleep(time.Duration(
			MinTTFTMs+rand.Intn(MaxTTFTMs-MinTTFTMs),
		) * time.Millisecond)

		lastMsg := ""
		if len(req.Messages) > 0 {
			lastMsg = req.Messages[len(req.Messages)-1].Content
		}

		response := generateMockResponse(lastMsg)
		words := strings.Split(response, " ")

		for _, w := range words {
			textChan <- w + " "
			time.Sleep(time.Duration(
				MinTokenDelay+rand.Intn(MaxTokenDelay-MinTokenDelay),
			) * time.Millisecond)
		}

		usage := estimateTokens(response)
		textChan <- fmt.Sprintf(
			"\n\n\033[90m[usage: %d tokens simulated]\033[0m",
			usage,
		)
	}()

	return textChan, errChan
}

////////////////////////////////////////////////////////////////////////////////
// NON-STREAM MODE
////////////////////////////////////////////////////////////////////////////////

func (c *MockLLMConnector) ChatCompletion(
	req ChatCompletionRequest,
) (string, UsageStats, error) {

	if err := simulateFailure(); err != nil {
		return "", UsageStats{}, err
	}

	lastMsg := ""
	for i := len(req.Messages) - 1; i >= 0; i-- {
		if req.Messages[i].Role == "user" {
			lastMsg = req.Messages[i].Content
			break
		}
	}

	resp := generateMockResponse(lastMsg)

	usage := UsageStats{
		PromptTokens:     estimateTokens(lastMsg),
		CompletionTokens: estimateTokens(resp),
	}
	usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens

	return resp, usage, nil
}

////////////////////////////////////////////////////////////////////////////////
// DEMO (WORKFLOW READY)
////////////////////////////////////////////////////////////////////////////////

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("      VASTAR MOCK LLM CONNECTOR - PRODUCTION FINAL")
	fmt.Println("      MODE: OFFLINE | CI | WORKFLOW DEV")
	fmt.Println(strings.Repeat("=", 80))

	llm := NewMockLLMConnector("")
	defer llm.Close()

	// STREAMING DEMO
	fmt.Println("\n▶ Streaming Demo")
	fmt.Println(strings.Repeat("-", 80))

	req := ChatCompletionRequest{
		Messages: []Message{
			{Role: "user", Content: "Apa itu IPC dan kenapa penting?"},
		},
		Stream: true,
	}

	fmt.Print("\033[34mUser:\033[0m ", req.Messages[0].Content, "\n")
	fmt.Print("\033[33mAI:\033[0m ")

	tChan, eChan := llm.ChatCompletionStream(req)
	for {
		select {
		case t, ok := <-tChan:
			if !ok {
				fmt.Println()
				goto NonStream
			}
			fmt.Print(t)
		case err := <-eChan:
			log.Fatalf("STREAM ERROR: %v", err)
		}
	}

NonStream:
	// NON-STREAM DEMO
	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("▶ Non-stream Completion Demo")

	resp, usage, err := llm.ChatCompletion(ChatCompletionRequest{
		Messages: []Message{
			{Role: "user", Content: "Apa itu Vastar?"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\033[33mAI:\033[0m", resp)
	fmt.Printf("\nUsage: %+v\n", usage)

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("✅ MOCK VALIDATION COMPLETE – READY FOR CI & WORKFLOWS")
}
