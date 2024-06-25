package interceptor

import (
	"context"

	"connectrpc.com/connect"
)

// AuthFunc is the function that handles the authentication logic
func AuthInterceptor(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		token := req.Header().Get("Authorization")

		// Here you would validate the token and extract the user ID
		// For example, decode a JWT token to get the user ID
		// For simplicity, let's assume the token is the user ID
		userID := token

		// Add the user ID to the context
		ctx = context.WithValue(ctx, "userID", userID)

		return next(ctx, req)
	}
}
