package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic" // Penting untuk akurasi metrik saat testing berat
	"time"

	// Menggunakan SDK asli [cite: 3535]
	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
	ipc "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang/protocol"
)

// Menggunakan atomic untuk keamanan data saat Stress Test [cite: 3490-3492]
var (
	requestCount  uint64
	latencySumMs  uint64 // dalam microseconds agar lebih presisi
)

func main() {
	// 1. Inisialisasi Client dengan Tenant-ID sesuai instruksi mentor
	client, err := vastar.NewRuntimeClientWithOptions(
		60*time.Second, 
		"rizka-tenant-001", // Isolasi tenant
		"benchmark-workspace",
	)
	if err != nil {
		log.Fatalf("❌ BOOT ERROR: %v", err)
	}
	defer client.Close()

	// 2. Metrics Handler (Standard Prometheus Format) [cite: 72, 3543]
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		count := atomic.LoadUint64(&requestCount)
		totalLat := atomic.LoadUint64(&latencySumMs)
		
		avgLat := 0.0
		if count > 0 {
			avgLat = float64(totalLat) / float64(count) / 1000.0 // Konversi ke ms
		}

		fmt.Fprintf(w, "# HELP vastar_processed_tasks_total Jumlah request terisolasi WASM\n")
		fmt.Fprintf(w, "vastar_processed_tasks_total %v\n", count)
		fmt.Fprintf(w, "# HELP vastar_ipc_latency_ms Rata-rata latensi IPC (FlatBuffers)\n")
		fmt.Fprintf(w, "vastar_ipc_latency_ms %v\n", avgLat)
	})

	// 3. Gateway Logic (Adapter Terkelola)
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := fmt.Sprintf("trace-%d", start.UnixNano())

		// Baca payload dari tester
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}

		// Simulasi Adapter Terkelola sesuai arsitektur PDF [cite: 1505-1514]
		// Kita hit 'httpbin' sebagai target mock external API
		vastarReq := vastar.POST("https://httpbin.org/post").
			WithBody(body).
			WithHeader("Content-Type", "application/json").
			WithHeader("X-Vastar-Trace", traceID).
			WithTraceID(traceID) // Distributed tracing 

		// Eksekusi via Runtime (Jalur IPC Aman) [cite: 2174, 3517]
		resp, err := client.ExecuteHTTP(vastarReq)
		
		duration := time.Since(start)

		if err != nil {
			// Kategorisasi error untuk laporan benchmarking [cite: 3544-3546]
			log.Printf("⚠️ [%s] IPC Failure: %v", traceID, err)
			http.Error(w, fmt.Sprintf("IPC Bridge Error: %v", err), 502)
			return
		}

		// Update metrik secara aman (atomic)
		atomic.AddUint64(&requestCount, 1)
		atomic.AddUint64(&latencySumMs, uint64(duration.Microseconds()))

		// Response ke Tester
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Runtime-Duration", fmt.Sprintf("%v", duration))
		
		// Menunjukkan status sukses sesuai protokol [cite: 3605]
		if resp.ErrorClass == ipc.ErrorClassSuccess {
			w.WriteHeader(int(resp.StatusCode))
			w.Write(resp.Body)
		} else {
			w.WriteHeader(500)
			fmt.Fprintf(w, `{"error": "Vastar Runtime Error", "class": "%v"}`, resp.ErrorClass)
		}
	})

	// Display Dashboard Info
	fmt.Println("--------------------------------------------------------")
	fmt.Println("🏆 VASTAR ENTERPRISE GATEWAY READY")
	fmt.Println("--------------------------------------------------------")
	fmt.Printf("🎯 Target Tester : http://localhost:3000/process\n")
	fmt.Printf("📊 Metrics Panel : http://localhost:3000/metrics\n")
	fmt.Printf("🛡️ Isolation Mode: Tenant-Based (WASM Ready)\n")
	fmt.Println("--------------------------------------------------------")

	log.Fatal(http.ListenAndServe(":3000", nil))
}