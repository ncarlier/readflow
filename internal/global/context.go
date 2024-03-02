package global

import (
	"context"

	"github.com/ncarlier/readflow/pkg/middleware"
)

type key int

const (
	// ContextUserID is the key used to store current user ID into the request context
	ContextUserID key = iota
	// ContextUser is the key used to store current user into the request context
	ContextUser
	// ContextIncomingWebhook is the key used to store incoming webhook into the request context
	ContextIncomingWebhook
	// ContextIsAdmin is true when the user is an admin
	ContextIsAdmin
	// ContextDownloader is the key used to store downloader service (uggly hack)
	ContextDownloader
)

// NewBackgroundContextWithValues create new context with same values than another context
func NewBackgroundContextWithValues(src context.Context) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextDownloader, src.Value(ContextDownloader))
	ctx = context.WithValue(ctx, middleware.ContextRequestID, src.Value(middleware.ContextRequestID))
	ctx = context.WithValue(ctx, ContextIncomingWebhook, src.Value(ContextIncomingWebhook))
	ctx = context.WithValue(ctx, ContextIsAdmin, src.Value(ContextIsAdmin))
	ctx = context.WithValue(ctx, ContextUser, src.Value(ContextUser))
	ctx = context.WithValue(ctx, ContextUserID, src.Value(ContextUserID))
	return ctx
}
