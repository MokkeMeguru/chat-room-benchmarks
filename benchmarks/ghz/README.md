# ghz test

## usage

wip

```
$ go run ./cmd/connect

$ ghz --config ./benchmarks/ghz/listen_only.json

$ grpcurl -proto api/proto-spec/chat/v1/chat.proto -plaintext -d '{"room_id": "room1"}' localhost:8080 chat.v1.ChatService.Receive

$ go run ./examples/connect-client
```
