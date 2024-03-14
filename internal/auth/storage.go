package auth

import (
	"context"
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=TokenStorage
type TokenStorage interface {
	SetToken(ctx context.Context, userId, token string, expire time.Duration) error
	GetToken(ctx context.Context, userId string) (string, error)
}
