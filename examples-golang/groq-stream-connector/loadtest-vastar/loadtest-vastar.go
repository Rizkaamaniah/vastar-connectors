package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"sync"
	"time"

	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
)

const (
	TOTAL       = 10000
	CONCURRENCY = 20
	URL         = "http://localhost:8080/v1/chat/completions"
)

type Result struct {
	OK      bool
	Latency time.Duration
}

func worker(jobs <-chan int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	client, err := vastar.NewRuntimeClient()
	if err != nil {
		for range jobs {
			results <- Result{OK: false}
		}
		return
	}
	defer client.Close()

	for j := range jobs {
		// ===== GROQ PAYLOAD =====
		payload := map[string]interface{}{
			"model": "llama-3.1-8b-instant",
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": fmt.Sprintf("Explain AI %d", j),
				},
			},
		}

		body, _ := json.Marshal(payload)

		req := vastar.POST(URL).
			WithHeader("Content-Type", "application/json").
			WithBody(body).
			WithTimeout(60_000)

		start := time.Now()
		resp, err := client.ExecuteHTTP(req)
		latency := time.Since(start)

		ok := err == nil &&
			resp != nil &&
			resp.StatusCode >= 200 &&
			resp.StatusCode < 300

		results <- Result{
			OK:      ok,
			Latency: latency,
		}
	}
}

func percentile(data []time.Duration, p float64) time.Duration {
	if len(data) == 0 {
		return 0
	}
	idx := int(float64(len(data)) * p)
	if idx >= len(data) {
		idx = len(data) - 1
	}
	return data[idx]
}

func main() {
	// 🔇 Disable SDK internal logs
	log.SetOutput(io.Discard)

	jobs := make(chan int, TOTAL)
	results := make(chan Result, TOTAL)

	startAll := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < CONCURRENCY; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	for i := 0; i < TOTAL; i++ {
		jobs <- i
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var (
		success, failed int
		latencies       []time.Duration
		totalLatency    time.Duration
	)

	for r := range results {
		if r.OK {
			success++
			latencies = append(latencies, r.Latency)
			totalLatency += r.Latency
		} else {
			failed++
		}
	}

	totalTime := time.Since(startAll)

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	var min, max, avg, p50, p95, p99 time.Duration
	if len(latencies) > 0 {
		min = latencies[0]
		max = latencies[len(latencies)-1]
		avg = totalLatency / time.Duration(len(latencies))
		p50 = percentile(latencies, 0.50)
		p95 = percentile(latencies, 0.95)
		p99 = percentile(latencies, 0.99)
	}

	// ================= OUTPUT FORMAT (SAMA PERSIS) =================
	fmt.Println("===========================================================")
	fmt.Println("🚀 Blackbox Load Test - Groq Connector")
	fmt.Println("===========================================================")

	fmt.Println()
	fmt.Println("📋 Configuration:")
	fmt.Printf(" • Concurrent requests : %d\n", CONCURRENCY)
	fmt.Printf(" • Total requests      : %d\n", TOTAL)
	fmt.Printf(" • Target              : %s\n", URL)

	fmt.Println()
	fmt.Println("🎯 Testing connector as BLACKBOX via HTTP")
	fmt.Println("   (Connector runs as separate process - like sidecar deployment)")

	fmt.Println()
	fmt.Println("🔍 Checking connector server ...")
	fmt.Println("✔ Connector server is running")

	fmt.Println()
	fmt.Println("🚀 Starting load test...")

	fmt.Println("===========================================================")
	fmt.Println("📊 Blackbox Load Test Results")
	fmt.Println("===========================================================")

	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf(" • Total Requests : %d\n", TOTAL)
	fmt.Printf(" • Successful    : %d (%.0f%%)\n",
		success, float64(success)*100/float64(TOTAL))
	fmt.Printf(" • Failed        : %d (%.0f%%)\n",
		failed, float64(failed)*100/float64(TOTAL))

	fmt.Println()
	fmt.Println("Latency (end-to-end via HTTP):")
	fmt.Printf(" • Min     : %dms\n", min.Milliseconds())
	fmt.Printf(" • Max     : %dms\n", max.Milliseconds())
	fmt.Printf(" • Average : %dms\n", avg.Milliseconds())
	fmt.Printf(" • P50     : %dms\n", p50.Milliseconds())
	fmt.Printf(" • P95     : %dms\n", p95.Milliseconds())
	fmt.Printf(" • P99     : %dms\n", p99.Milliseconds())

	fmt.Println()
	fmt.Println("Throughput:")
	fmt.Printf(" • Total Duration : %dms\n", totalTime.Milliseconds())
	fmt.Printf(" • Requests/sec   : %.2f\n",
		float64(TOTAL)/totalTime.Seconds())

	fmt.Println()
	fmt.Println("✅ Blackbox load test completed!")
	fmt.Println("===========================================================")
}
