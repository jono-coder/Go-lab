docker run --name go-lab \
 -p 8283:8282 \
 -e APP_HOST=localhost \
 -e APP_PORT=8282 \
 -e APP_ROOT=/lab \
 -e ENV=dev \
 -e DB_DRIVER=sqlite \
 -e DB_DSN="file:golab.db?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)" \
 -e AUTH_TOKEN_URL="https://dummy_url.x" \
 -e AUTH_CLIENT_ID="my_client_id" \
 -e AUTH_CLIENT_SECRET="my_client_secret" \
 golab:latest
