module backend-grpc-server

go 1.21

require (
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
	nhooyr.io/websocket v1.8.10 // WebSocket library
)

require github.com/google/uuid v1.3.1

require (
	github.com/cenkalti/backoff/v4 v4.1.1 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/rs/cors v1.7.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
)
