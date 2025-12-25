package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
)

// =============================================================================
// CONFIGURATION
// =============================================================================

const (
	ConnectorVersion = "2.2.0-Production"
	DefaultTimeoutMS = 10000
	MaxParallelLimit = 10
	MaxRetries       = 2
)

// =============================================================================
// DATA MODELS
// =============================================================================

type HealthStatus struct {
	Service    string `json:"service"`
	Status     string `json:"status"` // UP | DEGRADED | DOWN
	LatencyMS  int64  `json:"latency_ms"`
	Endpoint   string `json:"endpoint"`
	Error      string `json:"error,omitempty"`
	ErrorClass string `json:"error_class,omitempty"`
	CheckedAt  string `json:"checked_at"`
}

type SystemReport struct {
	ReportID      string         `json:"report_id"`
	OverallStatus string         `json:"overall_status"`
	SDKVersion    string         `json:"sdk_version"`
	Platform      string         `json:"platform"`
	Timestamp     string         `json:"timestamp"`
	ExecutionTime string         `json:"execution_time"`
	Checks        []HealthStatus `json:"checks"`
}

// =============================================================================
// CONNECTOR ENGINE
// =============================================================================

type HealthConnector struct {
	client *vastar.RuntimeClient
	logger *log.Logger
	sem    chan struct{}
}

func NewHealthConnector() (*HealthConnector, error) {
	client, err := vastar.NewRuntimeClient()
	if err != nil {
		return nil, fmt.Errorf("runtime unreachable: %w", err)
	}

	return &HealthConnector{
		client: client,
		logger: log.New(os.Stdout, "[HEALTHCHECK] ", log.LstdFlags),
		sem:    make(chan struct{}, MaxParallelLimit),
	}, nil
}

func (c *HealthConnector) Close() error {
	c.logger.Println("closing connector")
	return c.client.Close()
}

// =============================================================================
// CORE HEALTH LOGIC
// =============================================================================

func (c *HealthConnector) CheckRuntime() error {
	req := vastar.GET("http://localhost/health").WithTimeout(2000)
	_, err := c.client.ExecuteHTTP(req)
	return err
}

func (c *HealthConnector) CheckService(name, url string) HealthStatus {
	c.sem <- struct{}{}
	defer func() { <-c.sem }()

	start := time.Now()
	status := "UP"
	errMsg := ""
	errClass := "SUCCESS"

	req := vastar.GET(url).
		WithHeader("User-Agent", "Vastar-Healthcheck/2.2").
		WithTimeout(DefaultTimeoutMS)

	var resp *vastar.HTTPResponse
	var err error

	for i := 0; i <= MaxRetries; i++ {
		resp, err = c.client.ExecuteHTTP(req)
		if err == nil && resp != nil && resp.StatusCode < 500 {
			break
		}
		time.Sleep(time.Duration(100*(i+1)) * time.Millisecond)
	}

	latency := time.Since(start).Milliseconds()

	if err != nil {
		status = "DOWN"
		errMsg = err.Error()
		errClass = "NETWORK_ERROR"
	} else if resp.StatusCode >= 500 {
		status = "DEGRADED"
		errMsg = fmt.Sprintf("HTTP %d", resp.StatusCode)
		errClass = "REMOTE_FAILURE"
	}

	return HealthStatus{
		Service:    name,
		Status:     status,
		LatencyMS:  latency,
		Endpoint:   url,
		Error:      errMsg,
		ErrorClass: errClass,
		CheckedAt:  time.Now().Format(time.RFC3339),
	}
}

func (c *HealthConnector) Audit(services map[string]string) SystemReport {
	start := time.Now()
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := []HealthStatus{}
	overall := "HEALTHY"

	for name, url := range services {
		wg.Add(1)
		go func(n, u string) {
			defer wg.Done()
			res := c.CheckService(n, u)

			mu.Lock()
			results = append(results, res)
			if res.Status == "DOWN" {
				overall = "CRITICAL"
			} else if res.Status == "DEGRADED" && overall == "HEALTHY" {
				overall = "WARNING"
			}
			mu.Unlock()
		}(name, url)
	}

	wg.Wait()

	return SystemReport{
		ReportID:      fmt.Sprintf("RPT-%d", time.Now().Unix()),
		OverallStatus: overall,
		SDKVersion:    ConnectorVersion,
		Platform:      fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Timestamp:     time.Now().Format(time.RFC3339),
		ExecutionTime: time.Since(start).String(),
		Checks:        results,
	}
}

// =============================================================================
// MAIN
// =============================================================================

func main() {
	services := map[string]string{
		"OpenAI API":   "https://api.openai.com/v1/models",
		"Gemini API":   "https://generativelanguage.googleapis.com",
		"AI Simulator": "http://localhost:4545/health",
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("VASTAR HEALTHCHECK CONNECTOR - PRODUCTION")
	fmt.Println(strings.Repeat("=", 80))

	connector, err := NewHealthConnector()
	if err != nil {
		log.Fatal(err)
	}
	defer connector.Close()

	if err := connector.CheckRuntime(); err != nil {
		log.Fatal("runtime check failed:", err)
	}

	report := connector.Audit(services)
	out, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(out))

	if report.OverallStatus == "CRITICAL" {
		os.Exit(2)
	}
}
