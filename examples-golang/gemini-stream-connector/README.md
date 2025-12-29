Gemini Stream Connector (Vastar WF SDK)

Gemini Stream Connector adalah contoh implementasi streaming text generation menggunakan Google Gemini API dengan bantuan Vastar Workflow Connector SDK.

Connector ini mendukung dua mode eksekusi:

Simulator Mode (OpenAI-compatible, lokal)

Real Gemini API Mode (streamGenerateContent)

Project ini juga menyertakan script load test (k6) untuk mengukur performa endpoint simulator.

✨ Features

✅ Streaming response (Server-Sent Events)
✅ Dual mode: Simulator & Real Gemini API
✅ Abstraksi HTTP via Vastar Runtime Client
✅ OpenAI-compatible simulator fallback
✅ Cocok untuk workflow benchmarking & load test

📁 Project Structure
gemini-stream-connector/
├── main.go
├── gemini-loadtest.js
└── README.md

⚙️ Requirements

Go ≥ 1.21

k6

Vastar Workflow Connector SDK

(Optional) Gemini API Key

🚀 Running the Connector
1️⃣ Simulator Mode (Default)

Mode ini tidak membutuhkan API key dan menggunakan OpenAI-compatible simulator.

Environment
unset GEMINI_API_KEY
export GEMINI_BASE_URL=http://localhost:8080

Run
go run main.go


Expected output:

🤖 Gemini Stream Connector
Mode  : OPENAI SIMULATOR
Status: Simulator connected ✅

2️⃣ Real Gemini API Mode

Mode ini menggunakan endpoint resmi Google Gemini:

https://generativelanguage.googleapis.com

Environment
export GEMINI_API_KEY=your_api_key_here
unset GEMINI_BASE_URL

Run
go run main.go


Expected output:

Mode  : REAL GEMINI API
Base  : https://generativelanguage.googleapis.com

🧠 How It Works
🔹 Real Gemini API Mode

Menggunakan endpoint:

POST /v1beta/models/{model}:streamGenerateContent


Response dikirim sebagai event stream

Setiap potongan teks (part.text) diproses dan dikirim ke channel Go

🔹 Simulator Mode

Request dikonversi ke format OpenAI:

POST /v1/chat/completions


Streaming menggunakan format SSE (data: {...})

Cocok untuk local testing & load test tanpa API limit

🧪 Example Request (Go)
req := GeminiRequest{
    Contents: []Content{
        {
            Role: "user",
            Parts: []Part{
                {Text: "Explain quantum computing in simple terms."},
            },
        },
    },
}

📊 Load Testing with k6 (Simulator)

⚠️ Catatan:
Load test dilakukan ke simulator, bukan ke real Gemini API, untuk menghindari rate limit & biaya API.

File: gemini-loadtest.js
import http from 'k6/http';
import { check } from 'k6';

export const options = {
    vus: 50,
    duration: '30s',
};

export default function () {
    const url = 'http://localhost:8080/v1/chat/completions';

    const payload = JSON.stringify({
        model: "gpt-4",
        messages: [
            { role: "user", content: "Explain concurrency briefly" }
        ],
        stream: false
    });

    const res = http.post(url, payload, {
        headers: { 'Content-Type': 'application/json' },
        timeout: '60s',
    });

    check(res, {
        'status 200': r => r.status === 200,
    });
}

Run Load Test
k6 run gemini-loadtest.js

📈 Interpreting Results

Latency → diukur dari http_req_duration

Throughput → http_reqs / second

Error rate → harus ≈ 0% di simulator

Cocok sebagai baseline performance sebelum dibandingkan dengan:

workflow engine (Temporal / Vastar WF)

connector lain (Groq, OpenAI, dll)

🎯 Use Cases

Benchmark streaming vs non-streaming

Load test AI workflow tanpa API cost

Validasi integrasi Gemini API

Perbandingan simulator vs real API behavior

📌 Summary

Gemini Stream Connector ini berfungsi sebagai:

Reference implementation Gemini streaming

Tool untuk local testing & load testing

Baseline connector sebelum masuk workflow orchestration

hasil loadtest simulator gemini(latest):
k6 run gemini-loadtest.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: gemini-loadtest.js
        output: -

     scenarios: (100.00%) 1 scenario, 50 max VUs, 1m0s max duration (incl. graceful stop):
              * default: 50 looping VUs for 30s (gracefulStop: 30s)


     ✓ status 200

     checks.........................: 100.00% ✓ 6754       ✗ 0
     data_received..................: 403 MB  13 MB/s
     data_sent......................: 1.7 MB  57 kB/s
     http_req_blocked...............: avg=54.56µs  min=1.37µs  med=5.59µs   max=24.07ms  p(90)=9.03µs   p(95)=20.85µs
     http_req_connecting............: avg=32.15µs  min=0s      med=0s       max=11.61ms  p(90)=0s       p(95)=0s
     http_req_duration..............: avg=223.6ms  min=19.89ms med=217.82ms max=668.7ms  p(90)=322.94ms p(95)=363.62ms
       { expected_response:true }...: avg=223.6ms  min=19.89ms med=217.82ms max=668.7ms  p(90)=322.94ms p(95)=363.62ms
     http_req_failed................: 0.00%   ✓ 0          ✗ 6754
     http_req_receiving.............: avg=160.63ms min=84.33µs med=155.29ms max=585ms    p(90)=248.25ms p(95)=282.84ms
     http_req_sending...............: avg=151.31µs min=9.42µs  med=30.79µs  max=88.71ms  p(90)=97.09µs  p(95)=232.11µs
     http_req_tls_handshaking.......: avg=0s       min=0s      med=0s       max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=62.81ms  min=7.68ms  med=56.7ms   max=408.73ms p(90)=99ms     p(95)=118.09ms
     http_reqs......................: 6754    222.272846/s
     iteration_duration.............: avg=222.4ms  min=23.06ms med=217.82ms max=603.8ms  p(90)=320.44ms p(95)=355.95ms
     iterations.....................: 6754    222.272846/s
     vus............................: 50      min=50       max=50
     vus_max........................: 50      min=50       max=50


running (0m30.4s), 00/50 VUs, 6754 complete and 0 interrupted iterations
default ✓ [======================================] 50 VUs  30s


k6 run gemini-loadtest.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: gemini-loadtest.js
        output: -

     scenarios: (100.00%) 1 scenario, 100 max VUs, 1m0s max duration (incl. graceful stop):
              * default: 100 looping VUs for 30s (gracefulStop: 30s)


     ✓ status 200

     checks.........................: 100.00% ✓ 6790       ✗ 0
     data_received..................: 407 MB  13 MB/s
     data_sent......................: 1.8 MB  59 kB/s
     http_req_blocked...............: avg=138.3µs  min=1.68µs  med=5.2µs    max=28.92ms  p(90)=8.84µs   p(95)=20.24µs
     http_req_connecting............: avg=105.75µs min=0s      med=0s       max=28.81ms  p(90)=0s       p(95)=0s
     http_req_duration..............: avg=444.19ms min=99.16ms med=434.17ms max=1.15s    p(90)=620.37ms p(95)=683.32ms
       { expected_response:true }...: avg=444.19ms min=99.16ms med=434.17ms max=1.15s    p(90)=620.37ms p(95)=683.32ms
     http_req_failed................: 0.00%   ✓ 0          ✗ 6790
     http_req_receiving.............: avg=320.38ms min=17.86ms med=311.12ms max=912.15ms p(90)=480.33ms p(95)=534.35ms
     http_req_sending...............: avg=132.22µs min=9.35µs  med=29.74µs  max=48.99ms  p(90)=87.12µs  p(95)=205.91µs
     http_req_tls_handshaking.......: avg=0s       min=0s      med=0s       max=0s       p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=123.67ms min=6.37ms  med=113.84ms max=605.95ms p(90)=189.1ms  p(95)=221.68ms
     http_reqs......................: 6790    223.853947/s
     iteration_duration.............: avg=443.86ms min=99.51ms med=434.04ms max=1.16s    p(90)=620.04ms p(95)=683.62ms
     iterations.....................: 6790    223.853947/s
     vus............................: 100     min=100      max=100
     vus_max........................: 100     min=100      max=100


running (0m30.3s), 000/100 VUs, 6790 complete and 0 interrupted iterations
default ✓ [======================================] 100 VUs  30s


k6 run gemini-loadtest.js

          /\      |‾‾| /‾‾/   /‾‾/
     /\  /  \     |  |/  /   /  /
    /  \/    \    |     (   /   ‾‾\
   /          \   |  |\  \ |  (‾)  |
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: gemini-loadtest.js
        output: -

     scenarios: (100.00%) 1 scenario, 1500 max VUs, 2m30s max duration (incl. graceful stop):
              * api_load: 166.67 iterations/s for 2m0s (maxVUs: 500-1500, gracefulStop: 30s)


     ✓ status 200

     checks.........................: 100.00% ✓ 19346     ✗ 0
     data_received..................: 1.2 GB  9.4 MB/s
     data_sent......................: 5.1 MB  42 kB/s
     dropped_iterations.............: 655     5.324631/s
     http_req_blocked...............: avg=118.33µs min=2.19µs   med=7.36µs   max=55.56ms p(90)=18.91µs  p(95)=233.05µs
     http_req_connecting............: avg=86.6µs   min=0s       med=0s       max=55.38ms p(90)=0s       p(95)=0s
     http_req_duration..............: avg=1.7s     min=2.12ms   med=428.22ms max=9.64s   p(90)=5.55s    p(95)=6.42s
       { expected_response:true }...: avg=1.7s     min=2.12ms   med=428.22ms max=9.64s   p(90)=5.55s    p(95)=6.42s
     http_req_failed................: 0.00%   ✓ 0         ✗ 19346
     http_req_receiving.............: avg=1.19s    min=60.79µs  med=304.12ms max=7.96s   p(90)=3.78s    p(95)=4.63s
     http_req_sending...............: avg=168.42µs min=10.44µs  med=42.91µs  max=54.75ms p(90)=123.52µs p(95)=235.11µs
     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s      p(90)=0s       p(95)=0s
     http_req_waiting...............: avg=516.83ms min=823.14µs med=112.63ms max=3.22s   p(90)=1.79s    p(95)=2.04s
     http_reqs......................: 19346   157.26764/s
     iteration_duration.............: avg=1.7s     min=2.65ms   med=428.23ms max=9.64s   p(90)=5.54s    p(95)=6.4s
     iterations.....................: 19346   157.26764/s
     vus............................: 703     min=0       max=950
     vus_max........................: 962     min=500     max=962


running (2m03.0s), 0000/0962 VUs, 19346 complete and 0 interrupted iterations
api_load ✓ [======================================] 0000/0962 VUs  2m0s  166.67 iters/s