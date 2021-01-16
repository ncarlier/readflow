package service

import "errors"

// ErrCategoryNotFound if a category is not found
var ErrCategoryNotFound = errors.New("category not found")

// ErrIncomingWebhookNotFound if an incoming webhook service is not found
var ErrIncomingWebhookNotFound = errors.New("incoming webhook not found")

// ErrOutgoingWebhookNotFound if an outgoing webhook service is not found
var ErrOutgoingWebhookNotFound = errors.New("outgoing webhook not found")

// ErrDeviceNotFound if a device is not found
var ErrDeviceNotFound = errors.New("device not found")

// ErrUserQuotaReached if an user reach its quota
var ErrUserQuotaReached = errors.New("user quota reached")

// ErrOutgoingWebhookSend if an article can not be send to the outgoing webhook
var ErrOutgoingWebhookSend = errors.New("unable to send article to outgoing webhook")

// ErrArticleArchiving if an article can not be archived
var ErrArticleArchiving = errors.New("unable to archive article")
