import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '10s', target: 10 },
        { duration: '10s', target: 20 },
        { duration: '10s', target: 30 },
        { duration: '10s', target: 0 },
    ]
}

const TOKENS = [
    'user_one', 'user_two', 'user_three', 'user_four', 'user_five'
]

export default function() {
    const noToken = Math.random() < 0.3

    const headers = { 'Content-Type': 'application/json' }
    if (!noToken) {
        const token = TOKENS[Math.floor(Math.random() * TOKENS.length)]
        headers['Authorization'] = `Bearer ${token}`
    }

    const res = http.post('http://localhost/hello',
        JSON.stringify({ userId: '123', name: 'Nyella' }),
        { headers }
    )

    check(res, {
        'allowed': (r) => r.status === 200,
        'rate limited': (r) => r.status === 429,
    })
    sleep(1)
}