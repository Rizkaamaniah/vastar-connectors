package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
	ipc "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang/protocol"
)

// Menentukan struktur log yang kompleks [cite: 3374-3380]
type WorkflowLog struct {
	TenantID   string  `json:"tenant_id"`
	Activity   string  `json:"activity"`
	DurationMs float64 `json:"duration_ms"`
}

type DBCommand struct {
	SQL    string        `json:"sql"`
	Params []interface{} `json:"params"`
}

func main() {
	// 1. Inisialisasi Client dengan Tenant & Workspace ID (Standar Isolasi WASM)
	tenantID := "enterprise-tenant"
	client, err := vastar.NewRuntimeClientWithOptions(
		45*time.Second,
		tenantID,
		"production-env",
	)
	if err != nil {
		log.Fatalf("❌ CRITICAL ERROR: Gagal inisialisasi IPC Vastar: %v", err)
	}
	defer client.Close()

	// 2. Data yang akan diproses (Simulasi data dari Workflow Designer)
	logData := WorkflowLog{
		TenantID:   tenantID,
		Activity:   "Stress Test Database via Vastar IPC",
		DurationMs: 12.5,
	}

	// 3. Membangun Payload Query yang Kompleks [cite: 3511-3513]
	cmd := DBCommand{
		SQL:    "INSERT INTO workflow_logs (tenant_id, activity, duration_ms, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		Params: []interface{}{logData.TenantID, logData.Activity, logData.DurationMs, time.Now()},
	}
	payload, _ := json.Marshal(cmd)

	// 4. Spesifikasi Database URL (Sesuai arsitektur Runtime) [cite: 1511, 2322]
	dbURL := "postgres://postgres:admin123@localhost:5432/vastar_db"

	// 5. Eksekusi melalui Jalur Terkelola (Managed Adapter)
	req := vastar.POST(dbURL).
		WithBody(payload).
		WithHeader("X-Vastar-Operation", "transactional-insert").
		WithHeader("Content-Type", "application/json").
		WithTimeout(10000) // 10 detik timeout [cite: 3558]

	fmt.Println("🚀 Menjalankan Managed Database Adapter...")
	start := time.Now()

	resp, err := client.ExecuteHTTP(req)
	
	elapsed := time.Since(start)

	// 6. Penanganan Error yang Sangat Detail (Enterprise Standard) [cite: 3544-3546, 3571-3572]
	if err != nil {
		if connErr, ok := err.(*vastar.ConnectorError); ok {
			log.Printf("❌ Database Error [%s]: %s", connErr.ErrorClass, connErr.Message)
			if connErr.IsRetryable() {
				log.Println("🔄 Error ini bersifat sementara (Transient), aman untuk dicoba ulang.") 
			}
		}
		return
	}

	// 7. Verifikasi Status Sukses [cite: 3519-3520]
	if resp.ErrorClass == ipc.ErrorClassSuccess {
		fmt.Println("---------------------------------------------------------")
		fmt.Printf("✅ STATUS: SUCCESS (Code: %d)\n", resp.StatusCode)
		fmt.Printf("⏱️  LATENSI IPC: %v (FlatBuffers Optimized)\n", elapsed)
		fmt.Printf("📁 TENANT: %s\n", tenantID)
		fmt.Println("---------------------------------------------------------")
		fmt.Printf("📝 Response dari DB: %s\n", string(resp.Body))
	}
}