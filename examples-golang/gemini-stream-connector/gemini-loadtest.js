import http from 'k6/http';
import { check } from 'k6';

export const options = {
    scenarios: {
        api_load: {
            executor: 'constant-arrival-rate',
            rate: 10000, // 10.000 request
            timeUnit: '1m',
            duration: '2m',
            preAllocatedVUs: 500,
            maxVUs: 1500,
        },
    },
};


export default function() {
    const url = 'http://localhost:8080/v1/chat/completions';

    const payload = JSON.stringify({
        model: 'gemini-2.0-flash',
        messages: [
            { role: 'user', content: 'Explain concurrency briefly' }
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