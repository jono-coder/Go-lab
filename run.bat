docker run --rm ^
 -p 8283:8282 ^
 -e ENV=dev ^
 -e DB_DRIVER=sqlite ^
 -e DB_DSN="file:golab.db?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)" ^
 golab:latest
