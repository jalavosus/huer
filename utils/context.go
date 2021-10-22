package utils

import (
	"context"

	"github.com/jalavosus/huer/internal/config"
)

func TimeoutCtx() (context.Context, context.CancelFunc) {
	return TimeoutCtxFromCtx(context.Background())
}

func TimeoutCtxFromCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, config.DefaultContextTimeout)
}

func WithTimeoutCtx(f func(context.Context) error) error {
	ctx, cancel := TimeoutCtx()
	defer cancel()

	return f(ctx)
}