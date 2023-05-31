package helper

import (
	"context"
	"time"

	"github.com/ncarlier/readflow/pkg/constant"
)

// NewBackgroundContextWithValues create new context with same values than another context
func NewBackgroundContextWithValues(src context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	ctx = context.WithValue(ctx, constant.ContextDownloader, src.Value(constant.ContextDownloader))
	ctx = context.WithValue(ctx, constant.ContextIncomingWebhook, src.Value(constant.ContextIncomingWebhook))
	ctx = context.WithValue(ctx, constant.ContextIsAdmin, src.Value(constant.ContextIsAdmin))
	ctx = context.WithValue(ctx, constant.ContextRequestID, src.Value(constant.ContextRequestID))
	ctx = context.WithValue(ctx, constant.ContextUser, src.Value(constant.ContextUser))
	ctx = context.WithValue(ctx, constant.ContextUserID, src.Value(constant.ContextUserID))
	return ctx, cancel
}
