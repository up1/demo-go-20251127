# Demo go websocket + redis

## Start Redis
```
$docker compose up -d
$docker compose ps
```

## Start websocket server
```
$go mod tidy
$go run server.go
```

Access to chat ui
* http://localhost:8080/

## Start with PM2
```
$npm i -g pm2
$pm2 start pm2.config.js
$pm2 status
$pm2 log
```