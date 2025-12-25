📨 Vastar Telegram Notification Connector – Enterprise Edition

A production-ready notification connector for integrating Go-based workflow systems with the Telegram Bot API via the Vastar Connector SDK .
Designed for enterprise needs requiring real-time alerts , low latency , and high reliability using the FlatBuffers IPC pipeline .

📑 Project Summary
Item	Detail
Version	1.0.0-Stable
Status	Production Ready
Language	Go
Platform	Linux / WSL
Messaging API	Telegram Bot API
Transport	Unix Domain Sockets (UDS)
IPC Protocol	FlatBuffers (via Vastar SDK)
🏗️ Connector Function

This connector functions as a Notification Gateway in an enterprise workflow system, with the main use cases:

🔔 System Alerting

Send automatic notifications to Telegram when:

service down

error runtime

workflow failed

monitoring threshold exceeded

📊 Workflow & AI Reporting

Sending results:

output AI (Gemini, Groq, Claude)

workflow summary

automatic reports (cron / scheduler)

🛡️ Audit & Observability

All requests go through Vastar Runtime

Errors are classified using IPC ErrorClass

Safe from blocking / hanging processes (timeout enforced)

🧬 Data Path Architecture

The connector uses binary IPC paths to minimize direct HTTP overhead:

[ Go Application ]
        │
        │ FlatBuffers (IPC)
        ▼
[ Vastar Runtime ]
        │
        │ HTTPS (Managed)
        ▼
[ Telegram Bot API ]

Architectural Excellence

🚀 Low latency (sub-millisecond IPC)

🔒 Do not use net/httpdirectly

🧠 Timeout & lifecycle dikontrol runtime

🧩 Ready to integrate into workflow engine (n8n / scheduler)

🚀 How to Use
1️⃣ Telegram Bot Preparation

Open Telegram and search for @BotFather

Run the command/newbot

Save Bot Token

Get Chat ID :

Use @userinfobot(personal)

Or group ID (for group notification)

Make sure the bot has been started

2️⃣ Environment Configuration

In terminal (WSL / Linux):

export TELEGRAM_TOKEN='your_bot_token_here'
export TELEGRAM_CHAT_ID='your_chat_id_here'

3️⃣ Module Initialization

Go to the connector folder, then:

go mod init telegram-connector
go mod edit -replace github.com/fullstack-aidev/vastar-wf-connector-sdk-bin/sdk-golang=../../sdk-golang
go mod tidy

4️⃣ Running the Connector

Make sure Vastar Runtime is active , then run:

go run main.go

📊 Output Analysis
Output in Terminal
================================================================================
VASTAR TELEGRAM NOTIFICATION CONNECTOR (PRODUCTION)
================================================================================
[TELEGRAM] message sent successfully (latency: 320ms)
================================================================================

Output in Telegram App

Messages received by the bot:

🚀 Enterprise Status
Time: Mon, 02 Jan 2025 12:30:05 WIB
Status: ONLINE
Engine: Vastar IPC

🛡️ Production-Safe Characteristics

✔ Controlled timeout (anti-hang)
✔ IPC ErrorClass validation
✔ Structured logging
✔ No blocking network calls
✔ Ready to be used as a workflow action

🔗 Advanced Integration

This connector is ready to be combined with:

Gemini Stream Connector

Groq LPU Connector

Healthcheck Connector

Scheduler / Cron workflow

n8n Custom Node backend

📄 License & Ownership

© Vastar Technologies
Enterprise Messaging & Workflow Integration