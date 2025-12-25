🚀 Vastar Universal Durable Connector
📌 Overview Project

Vastar Universal Durable Connector is an Action-based Backend Logic implementation using Vastar Go SDK designed as a real connector service for n8n Custom Node and Vastar Workflow Designer .

Unlike conventional connectors, this system implements Durable Execution Architecture , where each Action is executed safely, isolated, and failure-resistant .
This Connector acts as a managed adapter , so the workflow (WASM) never communicates directly outside the system , but rather through a controlled Runtime Connector.

🧠 Core Architecture Concept

This connector follows Vastar's design principles:

Workflow Unit = WASM (Isolated)

External Communication = Connector Runtime

All failures are isolated per workflow

With this approach, the failure of one action will not damage other workflows or the host system .

🛠️ Key Features (Durable Engine)

This code is built on 3 main pillars of system resilience :

1️⃣ Exponential Backoff Retry

If the request fails, the system will retry with a gradual time delay:

1s → 2s → 4s → 8s → ... (maksimal 32s)


This mechanism maintains connection stability to:

external API

Microservice

Database HTTP-based

2️⃣ Circuit Breaker Pattern

If failure occurs again

CircuitThreshold = 5


For:

The circuit will open.

All new requests will be temporarily blocked (30 seconds)

The system is given time to recover.

This prevents:

API target overload

Infinite retry

Resource starvation

3️⃣ WASM Isolation Ready

This connector is ready to run in the Vastar Workflow environment because:

All requests viaRuntimeClient

There isn't anyhttp.DefaultClient

Supports tenant & workflow isolation

🏗️ Supported Actions

(n8n Dropdown Concept)

Action Name	Technical Logic	Use Case
Durable HTTP	HTTP GET / POST / PUT / DELETE dengan Retry & Circuit Breaker	External API integration
Durable HTTP Stream SSE	Persistent connection (SSE) with streaming buffer	Real-time stream / chatbot
Durable Webhook	Managed HTTP POST with retry	Callback & event trigger
Durable Health Check	Lightweight GET check	Monitoring service
Durable Delay	Controlled sleep execution	Workflow orchestration
Durable Fan-Out	Parallel HTTP execution (goroutine)	Multi-endpoint fetch
Durable File Download	HTTP GET + file persistence	Download asset / report
Durable Metrics Snapshot	Runtime statistics collector	Observability & monitoring
🔌 API Contract (Real Usage)
Endpoint
POST /execute

Request Format
{
  "action": "Durable HTTP",
  "method": "GET",
  "url": "https://httpbin.org/get",
  "retry": 2
}

Response Format
{
  "success": true,
  "data": { ... }
}

🚀 How to Run
1️⃣ Prerequisites

Vastar Connector Runtime is now active

The default runtime runs on:

127.0.0.1:5000

2️⃣ Installation
cd examples-golang/vastar-universal-connector
go mod tidy

3️⃣ Execution
go run main.go


The server will be active on:

http://localhost:8080

📊 Monitoring & Metrics

The connector provides runtime metrics snapshots :

requests→ Total incoming requests

success→ Request successful

retries→ Retry that occurs

circuit → Status circuit breaker

Contoh output:

{
  "requests": 25,
  "success": 23,
  "retries": 6,
  "circuit": false
}

🛡️ Security & Architecture Compliance

This connector complies with Vastar Architecture standards :

❌ Not usinghttp.DefaultClient

✅ All requests viaclient.ExecuteHTTP

✅ Binary IPC (FlatBuffers over UDS)

✅ Managed timeout & retry

✅ Failure isolation per workflow

With this design, the connector:

Safe to run multi-tenant

Does not harm the host system

Siap production & orchestration

🎯 Target Use Case

Backend logic untuk n8n Custom Node

Workflow orchestration di Vastar Workflow Designer

Secure integration ke:

external API

AI service

Streaming service

Webhook system

✅ Status

✔ Production-ready
✔ Mentor-compliant
✔ n8n-compatible
✔ Vastar Architecture aligned