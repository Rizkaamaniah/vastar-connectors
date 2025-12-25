🚀 Gemini Stream Connector – Enterprise Workflow Engine

Gemini Stream Connector is a production-ready connector that acts as an adapter between Go-based Workflows and the Google Gemini API , using the Vastar Connector SDK as the runtime IPC layer.

This connector is designed for enterprise workflow orchestration needs that require:

low latency,

strict timeout control,

and isolated execution via IPC (Inter-Process Communication).

📑 Project Summary

Connector Name : Gemini Stream Connector

Versi: 1.0.0 (Production Ready)

Model AI Default: Gemini 2.0 Flash

Runtime Engine: Vastar Workflow Runtime

Protokol IPC: FlatBuffers (Binary)

Transport: Unix Domain Socket (UDS)

Communication Method to AI : HTTP Streaming ( streamGenerateContent)

🏗️ System Architecture

The connector does not communicate directly vianet/http , but rather all AI requests pass through Vastar Runtime to ensure system security, isolation, and stability.

[Go Workflow App]
        |
        |  FlatBuffers (IPC)
        v
[Vastar Runtime]
        |
        |  HTTP Streaming
        v
[Google Gemini API]

Why is this architecture used?

Avoiding the HTTP client bottleneck directly

Prevent memory leaks & zombie goroutines

Safe to run in a multi-tenant environment

✨ Key Features
⚡ Low Overhead IPC

Internal communication using FlatBuffers

Lower latency than JSON over HTTP

🔄 Streaming Response (Chunk-based)

Using endpointsstreamGenerateContent

Gemini returns an array of stream chunks

Output is processed sequentially and securely

⏱️ Timeout Management

Timeout is controlled by Runtime ( RequestTimeoutMS)

Default: 120 seconds , safe for long prompts

🛡️ Enterprise Safety

All requests viaclient.ExecuteHTTP

Errors are classified viaipc.ErrorClass

Do not usehttp.DefaultClient

🧪 Workflow-Ready

Can be called directly from:

CLI

n8n Custom Node

Another HTTP/IPC based Workflow Engine

🛠️ Installation & Preparation
1️⃣ Prerequisites

Make sure:

Vastar Runtime is already running

Gemini API Key environment variable has been set

export GEMINI_API_KEY="your-gemini-api-key"

2️⃣ Dependency Synchronization

Di folder project:

go mod tidy

3️⃣ Running the Connector
go run main.go

🧠 Contoh Output Runtime
======================================================================
VASTAR GEMINI STREAM CONNECTOR (PRODUCTION)
Protocol : IPC / FlatBuffers
Runtime  : Managed
======================================================================
🤖 Gemini: Enterprises prefer IPC over HTTP because IPC eliminates network
serialization overhead, reduces latency, and allows tighter control
over execution boundaries in workflow systems.

Status: SUCCESS
======================================================================

📊 Audit & Observability

Each execution generates an implicit audit through:

IPC Error Class

Runtime Execution Status

Output Length Monitoring

Final status example:

Status: SUCCESS

🧩 Technical Details
Component	Mark
Model Default	gemini-2.0-flash
Timeout	120.000 ms
Streaming Format	Array of chunks
IPC Layer	FlatBuffers
Transport	Unix Domain Socket
JSON Parsing	Controlled & deterministic
🔐 Security Note (Important)

❌ Do not use net/httpdirectly

❌ Do not open public sockets

✅ All requests are monitored by Runtime

✅ Safe for multi-workflow execution

✅ Production Status

✔ Build clean
✔ No unused imports
✔ Runtime-managed
✔ Safe to use for demo / real workflow
✔ Ready to use for n8n / orchestration engine

Powered by
⚙️ Vastar Connector SDK
🤖 Google Gemini AI