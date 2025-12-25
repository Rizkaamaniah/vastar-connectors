🤖 Vastar Mock LLM Connector – Production Final (Offline)

A high-performance Large Language Model (LLM) Mock Connector designed for offline simulation in development workflows, CI/CD testing, and validation of AI systems without any external API or runtime dependencies .

This project simulates the behavior of modern LLM (streaming, token usage, latency, error rate) with a deterministic & production-safe approach making it very suitable for enterprise pipelines.

📑 Project Summary

Version : 2.5.0-Production

Status: Production-Ready (Offline Simulation)

Operation Mode :

CI / Unit Test

Workflow Development

Offline AI Simulation

Dependency: ✅ Go Standard Library Only

External API : ❌ None

Runtime / IPC : ❌ None

🏗️ Simulation Architecture (Offline)

This connector runs entirely inside a Go process without any network communication or IPC:

[ Go Application ]
        │
        ▼
[ Mock LLM Engine ]
  - Latency Simulation
  - Streaming Token
  - Token Usage Estimation
  - Error Injection


This architecture ensures stable, fast, and reproducible test results in a CI environment.

✨ Enterprise Features
🔁 Multi-turn Conversation

Simulates chat history for context and dialogue flow testing.

⏱️ Natural Streaming Simulation

Time-To-First-Token (TTFT)

Variable delay per token

Phased output like a real LLM

📊 Token Usage Tracking

Prompt Tokens

Completion Tokens

Total Tokens
Used to simulate future model cost estimation.

⚠️ Configurable Error Injection

Random error simulation (rate limit / transient failure)

Useful for resilience testing

🧪 CI Friendly

No internet needed

No API key needed

Not flakey

Deterministic

🚀 How to Run
1️⃣ Prerequisites

Go 1.21+

No need for additional runtimes, SDKs, or services

2️⃣ Module Initialization
go mod init vastar-mock-llm


No need go mod tidybecause there are no external dependencies

3️⃣ Execution
go run main.go

📊 Contoh Output
Output Terminal
================================================================================
      VASTAR MOCK LLM CONNECTOR - PRODUCTION FINAL
      MODE: OFFLINE | CI | WORKFLOW DEV
================================================================================

▶ Streaming Demo
--------------------------------------------------------------------------------
User: Apa itu IPC dan kenapa penting?
AI: IPC memungkinkan komunikasi antar proses dengan latensi rendah dibanding HTTP tradisional.

[usage: 36 tokens simulated]

--------------------------------------------------------------------------------
▶ Non-stream Completion Demo
AI: Vastar Workflow SDK adalah sistem IPC berbasis FlatBuffers yang dirancang untuk performa tinggi dan workflow enterprise.

Usage: {PromptTokens:8 CompletionTokens:16 TotalTokens:24}
================================================================================
✅ MOCK VALIDATION COMPLETE – READY FOR CI & WORKFLOWS

🛠️ Simulation Configuration (Code Level)

Simulation values ​​can be set directly in the code:

MinTTFTMs      = 300
MaxTTFTMs      = 800
MinTokenDelay  = 40
MaxTokenDelay  = 120
ErrorRatePct   = 5
MaxTokenFactor = 2

🧪 Use Case Enterprise

CI pipeline without API cost

AI workflow development

Load & resilience testing

Team training without production risk

Mock replacement for OpenAI / Groq / Claude

✅ Conclusion

Vastar Mock LLM Connector – Production Final is an AI simulation solution that:

Stable

Fast

Offline

Safe for enterprise

Ready to use directly in CI & workflow n

Use this mock in development, switch to the real connector in production — without changing the application architecture .