package connect

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"github.com/MokkeMeguru/chat-benchmarks/internal/infrastructure/connect/handler"
	"github.com/MokkeMeguru/chat-benchmarks/internal/infrastructure/connect/interceptor"
	"github.com/MokkeMeguru/chat-benchmarks/internal/infrastructure/connect/proto/chat/v1/chatv1connect"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Run(ctx context.Context, logger *slog.Logger) error {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, dbUser, dbName, dbPassword)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(150)
	db.SetConnMaxLifetime(10 * time.Second)

	// TODO fix redis spec for heavy use
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", redisHost, redisPort),
		PoolSize:     200,
		MinIdleConns: 150,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		PoolTimeout:  10 * time.Second,
	})

	mux := http.NewServeMux()
	chatService := handler.NewChatServiceHandler(db, rdb, logger)
	intercepters := connect.WithInterceptors(connect.UnaryInterceptorFunc(interceptor.AuthInterceptor))
	mux.Handle(chatv1connect.NewChatServiceHandler(chatService, intercepters))

	httpServer := &http.Server{
		Addr: ":8080",
		Handler: h2c.NewHandler(mux, &http2.Server{
			MaxConcurrentStreams: 6000,
		}),
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("Starting server on :8080")
		serverErrors <- httpServer.ListenAndServe()
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	select {
	case <-quit:
		logger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			logger.With("err", err).Error("Could not gracefully shut down the server")
			return err
		}

	case err := <-serverErrors:
		if err != nil {
			logger.With("err", err).Error("Server error", err)
			return err
		}
	}

	return nil
}
