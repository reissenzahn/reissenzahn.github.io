package main

import (
	"context"
	"fmt"
)

// the key must satisfy comparability; that is, the equality operators == and != need to return correct results when used

// values returned must be safe to access from multiple goroutines

// Use context values only for request-scoped data that transits processes and API boundaries, not for passing optional parameters to functions

type ctxKey int

const (
	ctxUserId ctxKey = iota
	ctxAuthToken
)

func UserId(ctx context.Context) string {
	return ctx.Value(ctxUserId).(string)
}

func AuthToken(ctx context.Context) string {
	return ctx.Value(ctxAuthToken).(string)
}

func ProcessRequest(userId, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserId, userId)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf("handling response for %v (%v)", UserId(ctx), AuthToken(ctx))
}

func main() {
	ProcessRequest("jane", "123")
}
