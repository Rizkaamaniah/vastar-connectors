package main

import (
	"io"
	"log"
	"net/http"
	"time"

	vastar "github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang"
)

const SIMULATOR_URL = "http://localhost:8080/v1/chat/completions"

func vastarHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	client, err := vastar.NewRuntimeClient()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer client.Close()

	req := vastar.POST(SIMULATOR_URL).
		WithHeader("Content-Type", "application/json").
		WithBody(body).
		WithTimeout(60_000)

	resp, err := client.ExecuteHTTP(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(resp.StatusCode))
	w.Write(resp.Body)

	log.Printf("[VASTAR] %v\n", time.Since(start))
}

func main() {
	http.HandleFunc("/vastar/chat", vastarHandler)

	log.Println("Vastar HTTP Service running on :9100")
	log.Fatal(http.ListenAndServe(":9100", nil))
}
