import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
    vus: 5, // kecilin dulu (API real)
    duration: "20s",
    thresholds: {
        http_req_duration: ["p(95)<2000"], // real API lebih lambat
        http_req_failed: ["rate<0.05"],
    },
};

const BASE_URL = "https://api.groq.com/openai/v1";
const API_KEY = __ENV.GROQ_API_KEY;

export default function() {
    const url = `${BASE_URL}/chat/completions`;

    const payload = JSON.stringify({
        model: "llama-3.1-8b-instant",
        messages: [{
            role: "user",
            content: "Explain concurrency in simple terms.",
        }, ],
        stream: false,
        max_tokens: 100,
        temperature: 0.7,
    });

    const headers = {
        "Content-Type": "application/json",
        Authorization: `Bearer ${API_KEY}`,
    };

    const res = http.post(url, payload, { headers });

    check(res, {
        "status is 200": (r) => r.status === 200,
        "latency < 2s": (r) => r.timings.duration < 2000,
    });

    sleep(1);
}