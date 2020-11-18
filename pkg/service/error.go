package service

import "errors"

// ErrCategoryNotFound if a category is not found
var ErrCategoryNotFound = errors.New("category not found")

// ErrInboundServiceNotFound if an inbound service is not found
var ErrInboundServiceNotFound = errors.New("inbound service not found")

// ErrOutboundServiceNotFound if an outbound service is not found
var ErrOutboundServiceNotFound = errors.New("outbound service not found")

// ErrDeviceNotFound if a device is not found
var ErrDeviceNotFound = errors.New("device not found")

// ErrUserQuotaReached if an user reach its quota
var ErrUserQuotaReached = errors.New("user quota reached")

// ErrOutboundServiceSend if an article can not be send to the outbound service
var ErrOutboundServiceSend = errors.New("unable to send article to outbound service")
