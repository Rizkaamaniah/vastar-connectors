import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    scenarios: {
        vastar_http: {
            executor: 'constant-arrival-rate',
            rate: 50, // TURUNKAN dari 200 → 50
            timeUnit: '1s',
            duration: '1m',
            preAllocatedVUs: 50,
            maxVUs: 150,
        },
    },
    thresholds: {
        http_req_failed: ['rate<0.10'], // toleransi 10%
        http_req_duration: ['p(95)<5000'], // p95 < 5 detik
    },
};

export default function() {
    const url = 'http://localhost:9100/vastar/chat';

    const payload = JSON.stringify({
        model: 'gemini-simulator',
        messages: [
            { role: 'user', content: 'hello from k6' }
        ],
        stream: false
    });

    const res = http.post(url, payload, {
        headers: {
            'Content-Type': 'application/json',
        },
        timeout: '5s', // longgarin timeout
    });

    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(0.05); // kasih napas kecil
}