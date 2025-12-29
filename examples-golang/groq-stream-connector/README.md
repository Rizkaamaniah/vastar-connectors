Groq Stream Connector (Vastar WF SDK)

Groq Stream Connector adalah contoh implementasi streaming chat completion berbasis OpenAI-compatible API menggunakan Vastar Workflow Connector SDK.
Connector ini mendukung dua mode:

Simulator Mode (local, tanpa API key)

Real Groq API Mode

Project ini juga dilengkapi dengan load testing menggunakan k6 untuk mengukur latency, throughput, dan error rate.

✨ Features

✅ Streaming Chat Completion (SSE / text/event-stream)

✅ OpenAI-compatible request & response

✅ Dual mode: Simulator & Real Groq API

✅ Built-in simulator connectivity test

✅ Load testing (simulator & real API)

✅ Cocok untuk baseline performance comparison

📁 Project Structure
groq-stream-connector/
├── main.go
├── groq-loadtest-simulator.js
├── groq-loadtest-real.js
└── README.md

⚙️ Requirements

Go ≥ 1.21

k6

Vastar Workflow Connector SDK

(Optional) Groq API Key

🚀 Running the Connector
1️⃣ Simulator Mode (Recommended for Development & Load Test)

Simulator tidak memerlukan API key.

Environment
unset GROQ_API_KEY
export GROQ_BASE_URL=http://localhost:4545

Run
go run main.go


Expected output:

🤖 Groq Stream Connector
Mode  : RAI SIMULATOR
Testing simulator connection...
Simulator OK ✅

2️⃣ Real Groq API Mode
Environment
export GROQ_API_KEY=your_api_key_here
unset GROQ_BASE_URL

Run
go run main.go


Expected output:

Mode  : REAL GROQ API
Base  : https://api.groq.com/openai

🧠 How It Works

Connector mengirim request ke endpoint:

POST /v1/chat/completions


Response diproses sebagai Server-Sent Events (SSE)

Setiap chunk streaming dikirim ke channel Go

Connector dapat digunakan langsung atau sebagai bagian workflow engine

📊 Load Testing with k6

Load test dilakukan tanpa streaming (stream: false) untuk hasil yang stabil.

🔹 Load Test – Simulator

File: groq-loadtest-simulator.js

import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
    vus: 5,
    duration: "20s",
    thresholds: {
        http_req_duration: ["p(95)<2000"],
        http_req_failed: ["rate<0.01"],
    },
};

const BASE_URL = "http://localhost:8080";

export default function () {
    const url = `${BASE_URL}/v1/chat/completions`;

    const payload = JSON.stringify({
        model: "llama-3.1-8b-instant",
        messages: [
            { role: "user", content: "Explain concurrency in simple terms." }
        ],
        stream: false,
        max_tokens: 100,
        temperature: 0.7,
    });

    const res = http.post(url, payload, {
        headers: { "Content-Type": "application/json" },
    });

    check(res, {
        "status is 200": (r) => r.status === 200,
        "latency < 2s": (r) => r.timings.duration < 2000,
    });

    sleep(1);
}


Run:

k6 run groq-loadtest-simulator.js

🔹 Load Test – Real Groq API

File: groq-loadtest-real.js

import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
    vus: 5,
    duration: "20s",
    thresholds: {
        http_req_duration: ["p(95)<2000"],
        http_req_failed: ["rate<0.05"],
    },
};

const BASE_URL = "https://api.groq.com/openai/v1";
const API_KEY = __ENV.GROQ_API_KEY;

export default function () {
    const url = `${BASE_URL}/chat/completions`;

    const payload = JSON.stringify({
        model: "llama-3.1-8b-instant",
        messages: [
            { role: "user", content: "Explain concurrency in simple terms." }
        ],
        stream: false,
        max_tokens: 100,
        temperature: 0.7,
    });

    const res = http.post(url, payload, {
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${API_KEY}`,
        },
    });

    check(res, {
        "status is 200": (r) => r.status === 200,
        "latency < 2s": (r) => r.timings.duration < 2000,
    });

    sleep(1);
}


Run:

GROQ_API_KEY=your_api_key_here k6 run groq-loadtest-real.js

📈 Notes on Results

Simulator Mode

Latency sangat rendah

Error rate ≈ 0%

Cocok untuk baseline & stress test

Real API Mode

Latency lebih tinggi

Bisa terkena rate limit

Cocok untuk real-world validation

🎯 Use Cases

Performance comparison (pure HTTP vs connector vs workflow)

Workflow engine benchmarking

Streaming vs non-streaming analysis

Load testing without API cost (simulator)

📌 Summary

Connector ini dirancang sebagai reference implementation untuk:

OpenAI-compatible streaming

Workflow-based AI execution

Load testing AI endpoints secara terstruktur


hasil loadtest dengan simulator(latest):
k6 run groq-loadtest-simulator.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: groq-loadtest-simulator.js
        output: -

     scenarios: (100.00%) 1 scenario, 5 max VUs, 50s max duration (incl. graceful stop):
              * default: 5 looping VUs for 20s (gracefulStop: 30s)


     ✓ status is 200
     ✓ latency < 2s

     checks.........................: 100.00% ✓ 200      ✗ 0
     data_received..................: 6.2 MB  303 kB/s
     data_sent......................: 31 kB   1.5 kB/s
     http_req_blocked...............: avg=46.3µs   min=4.86µs   med=9.06µs  max=1.65ms   p(90)=19.26µs p(95)=35.69µs
     http_req_connecting............: avg=32.77µs  min=0s       med=0s      max=1.51ms   p(90)=0s      p(95)=5.46µs
   ✓ http_req_duration..............: avg=18.41ms  min=3.04ms   med=11.48ms max=176.89ms p(90)=32.86ms p(95)=40.57ms
       { expected_response:true }...: avg=18.41ms  min=3.04ms   med=11.48ms max=176.89ms p(90)=32.86ms p(95)=40.57ms
   ✓ http_req_failed................: 0.00%   ✓ 0        ✗ 100
     http_req_receiving.............: avg=9.07ms   min=366.03µs med=5.67ms  max=153.24ms p(90)=17.45ms p(95)=20.71ms
     http_req_sending...............: avg=110.64µs min=22.62µs  med=73.49µs max=1.21ms   p(90)=194.5µs p(95)=303.98µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s      max=0s       p(90)=0s      p(95)=0s
     http_req_waiting...............: avg=9.23ms   min=1.27ms   med=4.77ms  max=169.8ms  p(90)=17.01ms p(95)=22.91ms
     http_reqs......................: 100     4.878741/s
     iteration_duration.............: avg=1.01s    min=1s       med=1.01s   max=1.04s    p(90)=1.03s   p(95)=1.04s
     iterations.....................: 100     4.878741/s
     vus............................: 5       min=5      max=5
     vus_max........................: 5       min=5      max=5


running (20.5s), 0/5 VUs, 100 complete and 0 interrupted iterations
default ✓ [======================================] 5 VUs  20s

k6 run groq-loadtest-simulator.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: groq-loadtest-simulator.js
        output: -

     scenarios: (100.00%) 1 scenario, 10 max VUs, 50s max duration (incl. graceful stop):
              * default: 10 looping VUs for 20s (gracefulStop: 30s)


     ✓ status is 200
     ✓ latency < 2s

     checks.........................: 100.00% ✓ 400      ✗ 0
     data_received..................: 12 MB   599 kB/s
     data_sent......................: 63 kB   3.1 kB/s
     http_req_blocked...............: avg=37.76µs min=2.78µs   med=7.79µs  max=1.28ms   p(90)=18.74µs  p(95)=141.28µs
     http_req_connecting............: avg=17.84µs min=0s       med=0s      max=1.07ms   p(90)=0s       p(95)=3.86µs
   ✓ http_req_duration..............: avg=15.57ms min=2.29ms   med=10.97ms max=84.89ms  p(90)=30.5ms   p(95)=44.67ms
       { expected_response:true }...: avg=15.57ms min=2.29ms   med=10.97ms max=84.89ms  p(90)=30.5ms   p(95)=44.67ms
   ✓ http_req_failed................: 0.00%   ✓ 0        ✗ 200
     http_req_receiving.............: avg=9.67ms  min=217.13µs med=6.14ms  max=66.13ms  p(90)=20.87ms  p(95)=32.14ms
     http_req_sending...............: avg=78.96µs min=17.5µs   med=61.39µs max=380.33µs p(90)=152.22µs p(95)=195.57µs
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=5.81ms  min=941.81µs med=2.82ms  max=66.24ms  p(90)=11.6ms   p(95)=22.43ms
     http_reqs......................: 200     9.816449/s
     iteration_duration.............: avg=1.01s   min=1s       med=1.01s   max=1.08s    p(90)=1.03s    p(95)=1.04s
     iterations.....................: 200     9.816449/s
     vus............................: 10      min=10     max=10
     vus_max........................: 10      min=10     max=10


running (20.4s), 00/10 VUs, 200 complete and 0 interrupted iterations
default ✓ [======================================] 10 VUs  20s


k6 run groq-loadtest-simulator.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: groq-loadtest-simulator.js
        output: -

     scenarios: (100.00%) 1 scenario, 50 max VUs, 50s max duration (incl. graceful stop):
              * default: 50 looping VUs for 20s (gracefulStop: 30s)


     ✓ status is 200
     ✓ latency < 2s

     checks.........................: 100.00% ✓ 2000      ✗ 0
     data_received..................: 60 MB   2.9 MB/s
     data_sent......................: 313 kB  15 kB/s
     http_req_blocked...............: avg=355.61µs min=2.18µs   med=6.85µs  max=11.79ms  p(90)=12.91µs  p(95)=233.71µs
     http_req_connecting............: avg=328.11µs min=0s       med=0s      max=11.19ms  p(90)=0s       p(95)=15.07µs
   ✓ http_req_duration..............: avg=27.9ms   min=1.9ms    med=6.82ms  max=327.95ms p(90)=62.52ms  p(95)=166.87ms
       { expected_response:true }...: avg=27.9ms   min=1.9ms    med=6.82ms  max=327.95ms p(90)=62.52ms  p(95)=166.87ms
   ✓ http_req_failed................: 0.00%   ✓ 0         ✗ 1000
     http_req_receiving.............: avg=18.04ms  min=135.65µs med=4.1ms   max=292.02ms p(90)=41.7ms   p(95)=125.91ms
     http_req_sending...............: avg=92.47µs  min=12.52µs  med=52.95µs max=7.14ms   p(90)=122.91µs p(95)=207.53µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s      max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=9.76ms   min=924.43µs med=1.98ms  max=292.32ms p(90)=15.81ms  p(95)=44.07ms
     http_reqs......................: 1000    47.677564/s
     iteration_duration.............: avg=1.02s    min=1s       med=1s      max=1.33s    p(90)=1.06s    p(95)=1.16s
     iterations.....................: 1000    47.677564/s
     vus............................: 50      min=50      max=50
     vus_max........................: 50      min=50      max=50


running (21.0s), 00/50 VUs, 1000 complete and 0 interrupted iterations
default ✓ [======================================] 50 VUs  20s


k6 run groq-loadtest-simulator.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: groq-loadtest-simulator.js
        output: -

     scenarios: (100.00%) 1 scenario, 100 max VUs, 50s max duration (incl. graceful stop):
              * default: 100 looping VUs for 20s (gracefulStop: 30s)


     ✓ status is 200
     ✓ latency < 2s

     checks.........................: 100.00% ✓ 3826     ✗ 0
     data_received..................: 114 MB  5.4 MB/s
     data_sent......................: 599 kB  28 kB/s
     http_req_blocked...............: avg=516.96µs min=2.39µs   med=6.66µs  max=79.72ms  p(90)=11.14µs  p(95)=352.46µs
     http_req_connecting............: avg=414.88µs min=0s       med=0s      max=70.33ms  p(90)=0s       p(95)=89.85µs
   ✓ http_req_duration..............: avg=73.3ms   min=2.11ms   med=12.37ms max=669.33ms p(90)=235.55ms p(95)=388.16ms
       { expected_response:true }...: avg=73.3ms   min=2.11ms   med=12.37ms max=669.33ms p(90)=235.55ms p(95)=388.16ms
   ✓ http_req_failed................: 0.00%   ✓ 0        ✗ 1913
     http_req_receiving.............: avg=47.77ms  min=111.71µs med=8.3ms   max=564.7ms  p(90)=164.36ms p(95)=258.53ms
     http_req_sending...............: avg=767.54µs min=11.43µs  med=41.43µs max=66.68ms  p(90)=119.84µs p(95)=368.74µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s      max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=24.76ms  min=957.87µs med=3.23ms  max=421.61ms p(90)=57.16ms  p(95)=163.03ms
     http_reqs......................: 1913    90.80877/s
     iteration_duration.............: avg=1.07s    min=1s       med=1.01s   max=1.67s    p(90)=1.23s    p(95)=1.38s
     iterations.....................: 1913    90.80877/s
     vus............................: 2       min=2      max=100
     vus_max........................: 100     min=100    max=100


running (21.1s), 000/100 VUs, 1913 complete and 0 interrupted iterations
default ✓ [======================================] 100 VUs  20s


rizka@DESKTOP-O00VMAS:/mnt/c/vastar-connectors/examples-golang/groq-stream-connector$ k6 run groq-loadtest-simulator.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: groq-loadtest-simulator.js
        output: -

     scenarios: (100.00%) 1 scenario, 800 max VUs, 2m30s max duration (incl. graceful stop):
              * api_10k: 167.00 iterations/s for 2m0s (maxVUs: 300-800, gracefulStop: 30s)


     ✓ status 200

     checks.........................: 100.00% ✓ 20041      ✗ 0
     data_received..................: 1.2 GB  10 MB/s
     data_sent......................: 6.1 MB  51 kB/s
     http_req_blocked...............: avg=43.28µs  min=2.32µs   med=6.97µs  max=50.73ms  p(90)=10.29µs  p(95)=26.47µs
     http_req_connecting............: avg=24.06µs  min=0s       med=0s      max=50.61ms  p(90)=0s       p(95)=0s
   ✓ http_req_duration..............: avg=247.23ms min=2.08ms   med=97.6ms  max=2.61s    p(90)=762.44ms p(95)=1.06s
       { expected_response:true }...: avg=247.23ms min=2.08ms   med=97.6ms  max=2.61s    p(90)=762.44ms p(95)=1.06s
   ✓ http_req_failed................: 0.00%   ✓ 0          ✗ 20041
     http_req_receiving.............: avg=182.81ms min=70.39µs  med=65.45ms max=2.18s    p(90)=566.77ms p(95)=799.98ms
     http_req_sending...............: avg=110.76µs min=9.74µs   med=40.47µs max=49.46ms  p(90)=111.09µs p(95)=207.77µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s      max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=64.31ms  min=757.66µs med=25.69ms max=757.02ms p(90)=190.63ms p(95)=274.49ms
     http_reqs......................: 20041   166.398006/s
     iteration_duration.............: avg=247.6ms  min=2.71ms   med=98.52ms max=2.61s    p(90)=761.69ms p(95)=1.06s
     iterations.....................: 20041   166.398006/s
     vus............................: 92      min=1        max=215
     vus_max........................: 300     min=300      max=300


running (2m00.4s), 000/300 VUs, 20041 complete and 0 interrupted iterations
api_10k ✓ [======================================] 000/300 VUs  2m0s  167.00 iters/s