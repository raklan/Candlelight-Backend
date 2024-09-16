module candlelight-ruleengine

go 1.22.0

replace candlelight-models => ../candlelight-models

replace candlelight-api => ../candlelight-api

require candlelight-models v0.0.0-00010101000000-000000000000

require (
	candlelight-api v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.11.5
	golang.org/x/crypto v0.25.0
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	golang.org/x/net v0.27.0 // indirect
)
