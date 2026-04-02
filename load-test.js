import http from 'k6/http';
import { check } from 'k6';

export let options = {
    stages: [
        { duration: '10s', target: 5000 },  // ramp to 5000 users
        { duration: '10s', target: 10000 }, // ramp to 10000 users
        { duration: '10s', target: 20000 }, // ramp to 10000 users
        { duration: '10s', target: 0 },     // ramp back down
    ]
}

const JWTS = [
    'abc1234567',
    'xyz9876543',
    'def4561230',
    'ghi7890123',
    'jkl3456789'
]

export default function() {
    const jwt = JWTS[Math.floor(Math.random() * JWTS.length)]

    const res = http.post('http://192.168.23.129:8080',
        JSON.stringify({ userId: '123', name: 'Nyella' }),
        { headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${jwt}`
        }}
    )

    check(res, {
        'allowed': (r) => r.status === 200,
        'rate limited': (r) => r.status === 429,
    })
}