# A distributed go-rate-limiter
A distributed rate limiter in GO using a fixed window algorithm, backed by Redis, with NGINX load balancing. 


<img width="1140" height="681" alt="image" src="https://github.com/user-attachments/assets/fc84b32c-79da-471c-ae89-b7e335eb5a0a" />




# How it works

Incoming requests hit NGINX which load balances across multiple rate limiter instances using round-robin. Each instance applies a fixed window algorithm a counter is incremented in Redis for every request, and if the counter exceeds the limit within the current window the request is rejected with a `429 Too Many Requests` response.

Because all instances share the same Redis store, the rate limit is enforced consistently regardless of which instance handles the request.

<img width="1372" height="1117" alt="image" src="https://github.com/user-attachments/assets/36174ff7-7b79-4784-82e6-a58e073d2404" />


### Token-based limiting
If a request includes an `Authorization: Bearer <token>` header, the token is used as the rate limit key. This means each user has their own independent counter and window.

### IP-based limiting
If no token is present, the rate limiter falls back to the client's IP address as the key using the `X-Forwarded-For` header set by NGINX. This covers unauthenticated traffic such as requests to get a token in the first place.

### Window expiry
Rate limit counters are stored in Redis with a TTL. The expiry is set on the first request in a new window, so the window starts from the moment of the first request rather than on a fixed schedule. Once the TTL expires the counter is reset and the client can make requests again.

<img width="1461" height="470" alt="image" src="https://github.com/user-attachments/assets/e1790d01-236c-4ad8-a1ac-f93866207d2d" />


# Configuration

All configuration is done via the `.env` file in the project root. A `.env.example` file is provided with all available variables and their defaults. Copy it to get started:
```bash
cp .env.example .env
```

| Variable | Default | Description |
|----------|---------|-------------|
| REDIS_ADDR | redis:6379 | Redis connection address |
| RATE_LIMIT | 100 | Max requests per window |
| WINDOW | 1m | Rate limit window duration |
| PORT | 8080 | Port the rate limiter listens on |

Example — 50 requests per 30 seconds:
```env
REDIS_ADDR=redis:6379
RATE_LIMIT=50
WINDOW=30s
PORT=8080
```


# Requirements
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [k6](https://k6.io/) (for load testing only)

That's it no Go, Redis, or NGINX installation needed locally. Everything runs inside Docker.


# How to run it 
The project is containerized with Docker all services (Go rate limiter, Redis, NGINX) run via Docker Compose. No local dependencies required.
Running Docker Compose will get everything up and running. You can scale the rate limiter to however many instances you want. Ideally this will be automated with autoscaling via Kubernetes or similar in the future.
```bash 
docker compose up --build --force-recreate --scale ratelimiter=3
```

Expected output:
<img width="1442" height="597" alt="image" src="https://github.com/user-attachments/assets/b1bbc3b1-5d36-48e7-a14c-e2407700ff2f" />

To verify all containers are running:

```bash
docker ps
```

<img width="1117" height="148" alt="image" src="https://github.com/user-attachments/assets/3ffadada-0710-46b5-b79d-c0efb293fefe" />



# Load Testing

Load testing is done with [k6](https://k6.io/). Install it before running the tests.

**Windows**
```bash
winget install k6 --source winget
```

**macOS**
```bash
brew install k6
```

**Linux**
```bash
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

Or refer to the [official k6 installation docs](https://grafana.com/docs/k6/latest/set-up/install-k6/) for other distributions.

The test script simulates a realistic environment with a ramp up of virtual users sending requests to the rate limiter. 30% of requests are sent without a token, which triggers IP-based rate limiting. The remaining 70% rotate across 5 user tokens, which triggers token-based rate limiting.

| Stage | Duration | Virtual Users |
|-------|----------|---------------|
| Ramp up | 10s | 10 |
| Ramp up | 10s | 20 |
| Sustained | 10s | 30 |
| Ramp down | 10s | 0 |

Run the test script with:
```bash
k6 run load test.js
```

A successful test will show both `allowed` and `rate limited` checks firing, confirming that both token-based and IP-based rate limiting are working correctly.

<img width="1025" height="796" alt="k6 test output" src="https://github.com/user-attachments/assets/7ee2be1d-b8c1-4dd0-8938-58487685b102" />

Container logs are available for each request, indicating whether it was allowed or rate limited. In the example below the rate limit is configured to 100 requests per minute.

<img width="643" height="442" alt="container logs" src="https://github.com/user-attachments/assets/ddeffca1-bc04-45ed-8da7-8e4051de58b1" />
