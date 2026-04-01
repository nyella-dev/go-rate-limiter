#!/bin/bash

URL="http://127.0.0.1:8080"
ENDPOINT="/hello"

BODY='{
    "userId": "123",
    "name": "Nyella"
}'

JWTS=(
    "$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 10)"
    "$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 10)"
    "$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 10)"
    "$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 10)"
    "$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c 10)"
)

FAKE_IPS=(
    "192.168.1.1"
    "192.168.1.2"
    "192.168.1.3"
    "192.168.1.4"
    "192.168.1.5"
)

# arguments
REQUESTS=${1:-1}
USE_JWT=${2:-true}
MODE=${3:-"jwt"}   # jwt, ip, or cycle

for i in $(seq 1 $REQUESTS); do
    JWT=${JWTS[$((($i - 1) % 5))]}
    IP=${FAKE_IPS[$((($i - 1) % 5))]}

    # cycle alternates between jwt and ip each request
    if [ "$MODE" = "cycle" ]; then
        if [ $(($i % 2)) -eq 0 ]; then
            MODE_THIS_REQUEST="ip"
        else
            MODE_THIS_REQUEST="jwt"
        fi
    else
        MODE_THIS_REQUEST=$MODE
    fi

    if [ "$USE_JWT" = "false" ]; then
        AUTH_HEADER="Authorization: nil"
    elif [ "$MODE_THIS_REQUEST" = "ip" ]; then
        AUTH_HEADER="Authorization: nil"
    else
        AUTH_HEADER="Authorization: Bearer $JWT"
    fi

    echo "Request $i — Mode: $MODE_THIS_REQUEST — JWT: $JWT — IP: $IP"
    curl -X POST "$URL$ENDPOINT" \
        -H "Content-Type: application/json" \
        -H "$AUTH_HEADER" \
        -H "X-Forwarded-For: $IP" \
        -d "$BODY"
    echo ""
done


# # cycle between JWT and IP each request
##bash request.sh 10 true cycle

# only use JWT
##bash request.sh 10 true jwt

# only use IP
##bash request.sh 10 true ip

# no JWT at all
##bash request.sh 10 false