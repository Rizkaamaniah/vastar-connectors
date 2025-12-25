/**
 * VASTAR GOOGLE SHEETS CONNECTOR - ENTERPRISE AUTOMATION
 * Version: 1.0.0-Stable
 * -----------------------------------------------------------------------------
 * Arsitektur:
 * [App Go] <--- IPC (Unix Socket) ---> [Vastar Runtime] <--- HTTPS ---> [Google API]
 * -----------------------------------------------------------------------------
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// Vastar SDK Core
	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
)

// Menggunakan huruf kecil untuk internal struct (unexported) sesuai standar Go
type sheetsValueRange struct {
	Values [][]interface{} `json:"values"`
}

type googleSheetsConnector struct {
	client        *vastar.RuntimeClient
	spreadsheetID string
	apiKey        string
}

// newSheetsConnector menginisialisasi client dengan handshake biner
func newSheetsConnector(id string, key string) (*googleSheetsConnector, error) {
	client, err := vastar.NewRuntimeClient()
	if err != nil {
		return nil, fmt.Errorf("vastar runtime connection refused: %w", err)
	}
	return &googleSheetsConnector{
		client:        client,
		spreadsheetID: id,
		apiKey:        key,
	}, nil
}

// appendRow menulis data ke baris baru menggunakan FlatBuffers IPC
func (c *googleSheetsConnector) appendRow(sheetName string, rowData []interface{}) {
	payload := sheetsValueRange{
		Values: [][]interface{}{rowData},
	}
	body, _ := json.Marshal(payload)

	// Endpoint Google Sheets API v4
	url := fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s!A1:append?valueInputOption=USER_ENTERED&key=%s", 
		c.spreadsheetID, sheetName, c.apiKey)

	// Routing request melalui Vastar Runtime
	httpReq := vastar.POST(url).
		WithHeader("Content-Type", "application/json").
		WithBody(body)

	start := time.Now()
	resp, err := c.client.ExecuteHTTP(httpReq)
	
	if err != nil {
		fmt.Printf("\n❌ Error pada jalur IPC: %v\n", err)
		return
	}

	if resp.StatusCode == 200 {
		fmt.Printf("✅ Data berhasil ditulis ke Sheets (Latensi: %v)\n", time.Since(start))
	} else {
		fmt.Printf("\n❌ Google API Error %d: %s\n", resp.StatusCode, string(resp.Body))
		fmt.Println("👉 Pastikan Spreadsheet sudah di-Share ke email Service Account Anda!")
	}
}

func main() {
	// Ambil konfigurasi dari Environment Variables
	apiKey := os.Getenv("GOOGLE_API_KEY")
	spreadsheetID := os.Getenv("SHEETS_ID") // Set ini di terminal: export SHEETS_ID='...'

	// TAMPILAN OUTPUT BERSIH (ENTERPRISE STYLE)
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("%54s\n", "VASTAR SDK - GOOGLE SHEETS AUTOMATION ENGINE")
	fmt.Println(strings.Repeat("=", 80))

	if apiKey == "" || spreadsheetID == "" {
		log.Fatal("❌ Error: GOOGLE_API_KEY atau SHEETS_ID belum diatur di environment")
	}

	connector, err := newSheetsConnector(spreadsheetID, apiKey)
	if err != nil {
		log.Fatalf("Initialization Failed: %v", err)
	}
	defer connector.client.Close()

	fmt.Println("🐧 Connected via Unix Socket: /tmp/vastar-connector-runtime.sock")
	fmt.Printf("📊 Target Spreadsheet ID: %s...\n", spreadsheetID[:10])
	
	// Data yang akan kita kirim sebagai log
	// Kolom: [Waktu, Status, Sumber, Pesan]
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	rowToLog := []interface{}{timestamp, "STABLE", "VASTAR-SDK", "Enterprise Automation Handshake Success"}

	fmt.Println("\nStep 1: Writing data packet to Google Cloud...")
	connector.appendRow("Sheet1", rowToLog)

	fmt.Println("\n" + strings.Repeat("-", 80))
	fmt.Println("✅ Workflow Task Completed.")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("[SYSTEM] %s Releasing resources...\n", time.Now().Format("15:04:05"))
}