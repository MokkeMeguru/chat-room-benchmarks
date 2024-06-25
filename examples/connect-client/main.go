package main

import (
	"context"
	"log/slog"
	"os"

	"net/http"

	"connectrpc.com/connect"
	chatv1 "github.com/MokkeMeguru/chat-benchmarks/internal/infrastructure/connect/proto/chat/v1"
	"github.com/MokkeMeguru/chat-benchmarks/internal/infrastructure/connect/proto/chat/v1/chatv1connect"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := chatv1connect.NewChatServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	// Test TempCreateUser RPC
	req := connect.NewRequest(&chatv1.TempCreateUserRequest{Name: "user"})
	resp, err := client.TempCreateUser(context.Background(), req)
	if err != nil {
		logger.With("error", err).Error("TempCreateUser failed")
		os.Exit(1)
	}
	logger.With("user", resp.Msg.User).Info("Created user")

	// Test Send RPC
	sendReq := connect.NewRequest(&chatv1.SendRequest{RoomId: "room1", Message: "Hello World"})
	sendReq.Header().Set("Authorization", resp.Msg.User.UserId)
	sendResp, err := client.Send(context.Background(), sendReq)
	if err != nil {
		logger.With("error", err).Error("Send failed")
		os.Exit(1)
	}
	logger.With("messageId", sendResp.Msg.MessageId).Info("Sent message ID")

	// Test Receive RPC
	recvReq := connect.NewRequest(&chatv1.ReceiveRequest{RoomId: "room1"})
	stream, err := client.Receive(context.Background(), recvReq)
	if err != nil {
		logger.With("error", err).Error("Receive failed")
		os.Exit(1)
	}
	for stream.Receive() {
		msg := stream.Msg()
		logger.With("message", msg.Message).Info("Received message")
	}
	if err := stream.Err(); err != nil {
		logger.With("error", err).Error("Receive failed")
		os.Exit(1)
	}
}
