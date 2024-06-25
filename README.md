❯ ghz -d '{"room_id": "room1"}' --keepalive=3s --connections=10 --concurrency=10 --duration=20s --timeout=20s --proto api/proto-spec/chat/v1/chat.proto --call chat.v1.ChatService.Receive localhost:8080

