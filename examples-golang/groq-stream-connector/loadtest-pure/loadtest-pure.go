package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TotalRequests = 10000
	Concurrency   = 20
	TargetURL     = "http://localhost:8080/v1/chat/completions"
)

type Result struct {
	OK      bool
	Latency time.Duration
}

func worker(
	client *http.Client,
	jobs <-chan int,
	results chan<- Result,
	wg *sync.WaitGroup,
	completed *int64,
) {
	defer wg.Done()

	for j := range jobs {
		// 🔹 GROQ PAYLOAD (ONLY DIFFERENCE)
		payload := map[string]interface{}{
			"model": "llama-3.1-8b-instant",
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": fmt.Sprintf("Explain AI %d", j),
				},
			},
			"stream": false,
		}

		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", TargetURL, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		start := time.Now()
		resp, err := client.Do(req)
		latency := time.Since(start)

		if err != nil {
			results <- Result{OK: false, Latency: latency}
			atomic.AddInt64(completed, 1)
			continue
		}

		resp.Body.Close()

		results <- Result{
			OK:      resp.StatusCode >= 200 && resp.StatusCode < 300,
			Latency: latency,
		}

		atomic.AddInt64(completed, 1)
	}
}

func percentile(data []time.Duration, p float64) time.Duration {
	if len(data) == 0 {
		return 0
	}
	index := int(float64(len(data)) * p)
	if index >= len(data) {
		index = len(data) - 1
	}
	return data[index]
}

func main() {
	fmt.Println("===========================================================")
	fmt.Println("🚀 Blackbox Load Test - PURE GO HTTP (GROQ)")
	fmt.Println("===========================================================")
	fmt.Println()
	fmt.Println("📋 Configuration:")
	fmt.Printf(" • Concurrent requests : %d\n", Concurrency)
	fmt.Printf(" • Total requests      : %d\n", TotalRequests)
	fmt.Printf(" • Target              : %s\n", TargetURL)
	fmt.Println()
	fmt.Println("🎯 Testing connector as BLACKBOX via HTTP")
	fmt.Println("   (Pure Go client, no SDK, no sidecar)")
	fmt.Println()

	// ================= CHECK SERVER =================
	fmt.Println("🔍 Checking target server ...")

	client := &http.Client{Timeout: 60 * time.Second}

	checkReq, err := http.NewRequest(
		"POST",
		TargetURL,
		bytes.NewBuffer([]byte(`{"ping":"ok"}`)),
	)
	if err != nil {
		fmt.Println("✖ Failed to create request")
		return
	}
	checkReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(checkReq)
	if err != nil {
		fmt.Println("✖ Target server is NOT reachable")
		return
	}
	resp.Body.Close()

	fmt.Println("✔ Target server is running")
	fmt.Println()

	// ================= START TEST =================
	fmt.Println("🚀 Starting load test...")

	jobs := make(chan int, TotalRequests)
	results := make(chan Result, TotalRequests)

	var completed int64
	startAll := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < Concurrency; i++ {
		wg.Add(1)
		go worker(client, jobs, results, &wg, &completed)
	}

	// 🔹 Mentor-style progress (SAMA PERSIS)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			done := atomic.LoadInt64(&completed)
			percent := float64(done) / float64(TotalRequests) * 100
			fmt.Printf("\rProgress: %d/%d (%.1f%%)", done, TotalRequests, percent)
			if done >= TotalRequests {
				fmt.Println()
				return
			}
		}
	}()

	for i := 0; i < TotalRequests; i++ {
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
		} else {
			failed++
		}
		latencies = append(latencies, r.Latency)
		totalLatency += r.Latency
	}

	totalDuration := time.Since(startAll)

	sort.Slice(latencies, func(i, j int) bool {
		return latencies[i] < latencies[j]
	})

	avg := totalLatency / time.Duration(len(latencies))
	min := latencies[0]
	max := latencies[len(latencies)-1]
	p50 := percentile(latencies, 0.50)
	p95 := percentile(latencies, 0.95)
	p99 := percentile(latencies, 0.99)

	// ================= RESULT =================
	fmt.Println("===========================================================")
	fmt.Println("📊 Blackbox Load Test Results")
	fmt.Println("===========================================================")
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf(" • Total Requests : %d\n", TotalRequests)
	fmt.Printf(" • Successful    : %d (%.0f%%)\n",
		success, float64(success)/float64(TotalRequests)*100)
	fmt.Printf(" • Failed        : %d (%.0f%%)\n",
		failed, float64(failed)/float64(TotalRequests)*100)
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
	fmt.Printf(" • Total Duration : %dms\n", totalDuration.Milliseconds())
	fmt.Printf(" • Requests/sec   : %.2f\n",
		float64(TotalRequests)/totalDuration.Seconds())
	fmt.Println()
	fmt.Println("✅ Blackbox load test completed!")
	fmt.Println("===========================================================")
}
