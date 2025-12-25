🚀 Vastar Groq Connector (Production-Safe)

A high-performance AI connector for Groq LPU Inference integration using the Vastar Connector SDK .
This connector is designed as an enterprise workflow backend with a focus on stability, low latency, and IPC security , using production-safe buffered stream execution .

📋 Overview

This connector allows Go applications to communicate with the Groq Cloud API via the Vastar Runtime .
Instead of streaming directly from HTTP to the application, this connector uses a Buffered Stream Execution approach , where:

Groq still sends responses in modestream=true

Vastar Runtime manages and stabilizes data flow

Go applications process stream results incrementally from IPC buffers.

This approach avoids the risk of dropped connections, zombie processes, and streaming errors at the OS level , making it suitable for enterprise use and workflow engines.

✨ Key Features

Groq LPU Optimized
Leverages Groq's inference performance for low latency and high throughput.

Buffered Stream Execution (Production-Safe)
Groq streams are processed safely through runtime buffers, not direct SSE.

FlatBuffers IPC
Inter-process communication using a high-performance binary protocol.

Enterprise Error Handling
Using the official error classification from Vastar IPC ( ErrorClass).

Workflow Ready
Suitable as backend logic for n8n custom node or other workflow engines.

🏗️ System Architecture
[ Go Application ]
        |
        |  FlatBuffers IPC
        v
[ Vastar Runtime ]
        |
        |  HTTPS (Groq stream=true, buffered)
        v
[ Groq Cloud API (LPU Inference) ]


📌 Important note:
This connector is NOT a real-time SSE pass-through .
Groq streams are consolidated by the runtime to ensure system stability.

🚀 Usage Guide
1️⃣ Prerequisites

Go 1.21 or later

Vastar Connector Runtime is now running

Groq API Key from https://console.groq.com

2️⃣ Environment Configuration

Make sure the runtime is active and the API Key is set:

./start_runtime.sh
export GROQ_API_KEY="gsk_xxxxxxx"

3️⃣ Running the Connector
go mod tidy
go run main.go

📊 Performance Characteristics

Very low IPC Latency
(depends on host & runtime)

Execution Model
Buffered stream (safe for long workflows)

Stability
Does not depend on long SSE connections at the application level

Memory Usage
is light and controlled by runtime

🔐 Security & Compliance Architecture

This connector:

Do not usehttp.DefaultClient

All requests viaclient.ExecuteHTTP

Timeouts are managed by the runtime

Do not open network connections directly from the application.

Complies with Vastar Connector Architecture Guidelines .

📄 License

Copyright © Vastar Technologies
Used for internal development, workflow engine, and enterprise integration.

➡️ Next Steps

Integration as n8n Custom Node (dropdown action)

Adding retry & circuit breaker

Making the connector a Durable AI Action