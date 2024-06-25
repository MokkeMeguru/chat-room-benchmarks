package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/MokkeMeguru/chat-benchmarks/internal/domain/model"
	chatv1 "github.com/MokkeMeguru/chat-benchmarks/internal/infrastructure/connect/proto/chat/v1"
	"github.com/MokkeMeguru/chat-benchmarks/internal/infrastructure/connect/proto/chat/v1/chatv1connect"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	ulid "github.com/oklog/ulid/v2"
)

type ChatServiceHandler struct {
	chatv1connect.UnimplementedChatServiceHandler

	DB     *sql.DB
	RDB    *redis.Client
	Logger *slog.Logger
}

func NewChatServiceHandler(db *sql.DB, rdb *redis.Client, logger *slog.Logger) *ChatServiceHandler {
	return &ChatServiceHandler{
		DB:     db,
		RDB:    rdb,
		Logger: logger,
	}
}

func (h *ChatServiceHandler) TempCreateUser(ctx context.Context, req *connect.Request[chatv1.TempCreateUserRequest]) (*connect.Response[chatv1.TempCreateUserResponse], error) {
	user := &model.User{
		UserID: NewUUID(),
		Name:   req.Msg.Name,
	}
	if err := user.Insert(ctx, h.DB); err != nil {
		h.Logger.With("error", err).Error("failed to insert user")
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&chatv1.TempCreateUserResponse{User: &chatv1.User{
		UserId: user.UserID,
		Name:   user.Name,
	}}), nil
}

func (h *ChatServiceHandler) Send(ctx context.Context, req *connect.Request[chatv1.SendRequest]) (*connect.Response[chatv1.SendResponse], error) {
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("user ID not found in context"))
	}
	messageID, err := NewSequentialUUID()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	message := &model.Message{
		MessageID: messageID,
		RoomID:    req.Msg.RoomId,
		Content:   req.Msg.Message,
		UserID:    userID,
	}
	h.Logger.With("message", message).Info("send message")
	if err := message.Insert(ctx, h.DB); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if err := h.RDB.Publish(fmt.Sprintf("chat:%s", req.Msg.RoomId), messageID).Err(); err != nil {
		slog.Error("failed to publish message", err)
	}
	return connect.NewResponse(&chatv1.SendResponse{MessageId: messageID}), nil
}

func (h *ChatServiceHandler) Receive(ctx context.Context, req *connect.Request[chatv1.ReceiveRequest], stream *connect.ServerStream[chatv1.ReceiveResponse]) error {
	sub := h.RDB.Subscribe(fmt.Sprintf("chat:%s", req.Msg.RoomId))
	defer sub.Close()

	ch := sub.Channel()

	for {
		select {
		case msg := <-ch:
			messageID := msg.Payload

			// message, err := model.MessageByMessageID(ctx, h.DB, messageID)
			// if err != nil {
			// 	h.Logger.With("error", err).Error("failed to get message")
			// 	continue
			// }
			// user, err := model.UserByUserID(ctx, h.DB, message.UserID)
			// if err != nil {
			// 	h.Logger.With("error", err).Error("failed to get user")
			// 	continue
			// }
			res := &chatv1.ReceiveResponse{
				Message: &chatv1.Message{
					MessageId: messageID,
					Message:   "temp",
					User: &chatv1.User{
						UserId: "temp",
						Name:   "temp",
					},
				},
			}
			if err := stream.Send(res); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// utils
var (
	monotonic *safeMonotonicReader
)

type safeMonotonicReader struct {
	mtx sync.Mutex
	ulid.MonotonicReader
}

func init() {
	monotonic = &safeMonotonicReader{
		MonotonicReader: ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0),
	}
}

func NewSequentialUUID() (string, error) {
	id, err := ulid.New(ulid.Timestamp(time.Now()), monotonic)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func NewUUID() string {
	return uuid.New().String()
}
