package service

import "errors"

var (
	// ErrCategoryNotFound if a category is not found
	ErrCategoryNotFound = errors.New("category not found")
	// ErrIncomingWebhookNotFound if an incoming webhook service is not found
	ErrIncomingWebhookNotFound = errors.New("incoming webhook not found")
	// ErrOutgoingWebhookNotFound if an outgoing webhook service is not found
	ErrOutgoingWebhookNotFound = errors.New("outgoing webhook not found")
	// ErrDeviceNotFound if a device is not found
	ErrDeviceNotFound = errors.New("device not found")
	// ErrUserQuotaReached if an user reach its quota
	ErrUserQuotaReached = errors.New("user quota reached")
	// ErrOutgoingWebhookSend if an article can not be send to the outgoing webhook
	ErrOutgoingWebhookSend = errors.New("unable to send article to outgoing webhook")
	// ErrArticleDownload if an article can not be downloaded
	ErrArticleDownload = errors.New("unable to download article")
)
