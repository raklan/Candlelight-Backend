module candlelight-api

go 1.22.0

replace candlelight-ruleengine => ../candlelight-ruleengine

replace candlelight-models => ../candlelight-models

require (
	candlelight-models v0.0.0-00010101000000-000000000000
	candlelight-ruleengine v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.5.1
	github.com/rs/cors v1.11.1
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)
