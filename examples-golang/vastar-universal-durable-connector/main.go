package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
	ipc "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang/protocol"
)

//
// ============================================================================
// CONFIG
// ============================================================================
//

const (
	MaxBackoffTime   = 32 * time.Second
	CircuitThreshold = 5
	ServerPort       = ":8080"
)

//
// ============================================================================
// ACTION NAMES (n8n Dropdown)
// ============================================================================
//

const (
	ActionHTTP     = "Durable HTTP"
	ActionSSE      = "Durable HTTP Stream SSE"
	ActionWebhook  = "Durable Webhook"
	ActionHealth   = "Durable Health Check"
	ActionDelay    = "Durable Delay"
	ActionFanOut   = "Durable Fan-Out"
	ActionDownload = "Durable File Download"
	ActionMetrics  = "Durable Metrics Snapshot"
)

//
// ============================================================================
// DATA MODELS
// ============================================================================
//

type ExecuteRequest struct {
	Action     string            `json:"action"`
	Method     string            `json:"method,omitempty"`
	URL        string            `json:"url,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       json.RawMessage   `json:"body,omitempty"`
	Retry      int               `json:"retry,omitempty"`
	DelaySec   int               `json:"delay_sec,omitempty"`
	URLs       []string          `json:"urls,omitempty"`
	OutputFile string            `json:"output_file,omitempty"`
}

type ExecuteResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ConnectorStats struct {
	Requests uint64
	Success  uint64
	Retries  uint64
	Failures int32
	Circuit  bool
	LastFail time.Time
}

//
// ============================================================================
// CONNECTOR CORE (OOP STYLE)
// ============================================================================
//

type Connector struct {
	client *vastar.RuntimeClient
	stats  *ConnectorStats
	mu     sync.RWMutex
}

func NewConnector() (*Connector, error) {
	client, err := vastar.NewRuntimeClient()
	if err != nil {
		return nil, err
	}

	log.Println("🚀 Connected to Vastar Runtime")

	return &Connector{
		client: client,
		stats:  &ConnectorStats{},
	}, nil
}

//
// ============================================================================
// DURABLE EXECUTION ENGINE
// ============================================================================
//

func (c *Connector) execute(
	ctx context.Context,
	req *vastar.HTTPRequest,
	retry int,
) (*vastar.HTTPResponse, error) {

	atomic.AddUint64(&c.stats.Requests, 1)

	c.mu.RLock()
	if c.stats.Circuit && time.Since(c.stats.LastFail) < 30*time.Second {
		c.mu.RUnlock()
		return nil, errors.New("🚨 circuit breaker open")
	}
	c.mu.RUnlock()

	var lastErr error

	for i := 0; i <= retry; i++ {
		resp, err := c.client.ExecuteHTTP(req)

		if err == nil && resp.ErrorClass == ipc.ErrorClassSuccess {
			atomic.AddUint64(&c.stats.Success, 1)
			atomic.StoreInt32(&c.stats.Failures, 0)
			return resp, nil
		}

		lastErr = err
		atomic.AddUint64(&c.stats.Retries, 1)

		backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
		if backoff > MaxBackoffTime {
			backoff = MaxBackoffTime
		}

		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if atomic.AddInt32(&c.stats.Failures, 1) >= CircuitThreshold {
		c.mu.Lock()
		c.stats.Circuit = true
		c.stats.LastFail = time.Now()
		c.mu.Unlock()
		log.Println("🔥 Circuit breaker TRIPPED")
	}

	return nil, lastErr
}

//
// ============================================================================
// ACTION HANDLERS
// ============================================================================
//

func (c *Connector) handleHTTP(r ExecuteRequest) (interface{}, error) {
	log.Printf("🌐 [%s] %s %s", r.Action, r.Method, r.URL)

	var req *vastar.HTTPRequest

	switch r.Method {
	case "POST":
		req = vastar.POST(r.URL).WithBody(r.Body)
	case "PUT":
		req = vastar.PUT(r.URL).WithBody(r.Body)
	case "DELETE":
		req = vastar.DELETE(r.URL)
	default:
		req = vastar.GET(r.URL)
	}

	for k, v := range r.Headers {
		req.WithHeader(k, v)
	}

	resp, err := c.execute(context.Background(), req, r.Retry)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(resp.Body), nil
}

func (c *Connector) handleSSE(r ExecuteRequest) (interface{}, error) {
	log.Printf("📡 [SSE] %s", r.URL)

	req := vastar.GET(r.URL).
		WithHeader("Accept", "text/event-stream").
		WithTimeout(600000)

	resp, err := c.execute(context.Background(), req, 0)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(resp.Body))
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func (c *Connector) handleWebhook(r ExecuteRequest) (interface{}, error) {
	log.Printf("📨 [Webhook] %s", r.URL)
	r.Method = "POST"
	_, err := c.handleHTTP(r)
	return "webhook delivered", err
}

func (c *Connector) handleHealth(r ExecuteRequest) (interface{}, error) {
	log.Printf("🩺 [Health] %s", r.URL)
	_, err := c.handleHTTP(ExecuteRequest{
		Method: "GET",
		URL:    r.URL,
	})
	return err == nil, nil
}

func (c *Connector) handleDelay(r ExecuteRequest) (interface{}, error) {
	log.Printf("⏳ Delay %d seconds", r.DelaySec)
	time.Sleep(time.Duration(r.DelaySec) * time.Second)
	return fmt.Sprintf("delayed %d seconds", r.DelaySec), nil
}

func (c *Connector) handleFanOut(r ExecuteRequest) (interface{}, error) {
	log.Printf("🧨 Fan-Out %d targets", len(r.URLs))

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := []string{}

	for _, u := range r.URLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if _, err := c.handleHTTP(ExecuteRequest{Method: "GET", URL: url}); err == nil {
				mu.Lock()
				results = append(results, url)
				mu.Unlock()
			}
		}(u)
	}
	wg.Wait()

	return results, nil
}

func (c *Connector) handleDownload(r ExecuteRequest) (interface{}, error) {
	log.Printf("📥 Download %s → %s", r.URL, r.OutputFile)

	data, err := c.handleHTTP(r)
	if err != nil {
		return nil, err
	}

	return "file saved", os.WriteFile(
		r.OutputFile,
		data.(json.RawMessage),
		0644,
	)
}

func (c *Connector) handleMetrics() (interface{}, error) {
	return map[string]interface{}{
		"requests": atomic.LoadUint64(&c.stats.Requests),
		"success":  atomic.LoadUint64(&c.stats.Success),
		"retries":  atomic.LoadUint64(&c.stats.Retries),
		"circuit":  c.stats.Circuit,
	}, nil
}

//
// ============================================================================
// DISPATCHER (n8n BRAIN)
// ============================================================================
//

func (c *Connector) dispatch(r ExecuteRequest) (interface{}, error) {
	switch r.Action {
	case ActionHTTP:
		return c.handleHTTP(r)
	case ActionSSE:
		return c.handleSSE(r)
	case ActionWebhook:
		return c.handleWebhook(r)
	case ActionHealth:
		return c.handleHealth(r)
	case ActionDelay:
		return c.handleDelay(r)
	case ActionFanOut:
		return c.handleFanOut(r)
	case ActionDownload:
		return c.handleDownload(r)
	case ActionMetrics:
		return c.handleMetrics()
	default:
		return nil, fmt.Errorf("unknown action: %s", r.Action)
	}
}

//
// ============================================================================
// HTTP SERVER (REAL n8n ENTRYPOINT)
// ============================================================================
//

func main() {
	connector, err := NewConnector()
	if err != nil {
		log.Fatal(err)
	}
	defer connector.client.Close()

	server := &http.Server{
		Addr:         ServerPort,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var req ExecuteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		data, err := connector.dispatch(req)
		resp := ExecuteResponse{Success: err == nil, Data: data}
		if err != nil {
			resp.Error = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("🚀 Vastar Connector listening on", ServerPort)
	log.Fatal(server.ListenAndServe())
}
