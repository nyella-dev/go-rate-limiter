import http from 'k6/http';
import { check } from 'k6';

export let options = {
    vus: 10,         // 10 virtual users
    duration: '30s'  // for 30 seconds
}

const JWTS = [
    'abc1234567',
    'xyz9876543',
    'def4561230',
    'ghi7890123',
    'jkl3456789'
]

export default function() {
    // cycle through JWTs
    const jwt = JWTS[Math.floor(Math.random() * JWTS.length)]

    const res = http.post('http://127.0.0.1:8080/hello',
        JSON.stringify({ userId: '123', name: 'Nyella' }),
        { headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${jwt}`
        }}
    )

    // check response
    check(res, {
        'allowed': (r) => r.status === 200,
        'rate limited': (r) => r.status === 429,
    })
}