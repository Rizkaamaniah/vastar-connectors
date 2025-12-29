import http from "k6/http";
import { check } from "k6";

export const options = {
    scenarios: {
        api_10k: {
            executor: "constant-arrival-rate",
            rate: 167, // request per detik
            timeUnit: "1s",
            duration: "2m",
            preAllocatedVUs: 300,
            maxVUs: 800,
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.01"],
        http_req_duration: ["p(95)<2000"],
    },
};

export default function() {
    const res = http.post(
        "http://localhost:8080/v1/chat/completions",
        JSON.stringify({
            model: "llama-3.1-8b-instant",
            messages: [{ role: "user", content: "Explain concurrency simply." }],
            stream: false,
            max_tokens: 100,
            temperature: 0.7,
        }), { headers: { "Content-Type": "application/json" } }
    );

    check(res, {
        "status 200": r => r.status === 200,
    });
}