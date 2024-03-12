package auth

import (
	"context"
	"time"
)

type TokenStorage interface {
	SetToken(ctx context.Context, userId, token string, expire time.Duration) error
}
