package service

import "errors"

// ErrCategoryNotFound if a category is not found
var ErrCategoryNotFound = errors.New("category not found")

// ErrAPIKeyNotFound if an API key is not found
var ErrAPIKeyNotFound = errors.New("API key not found")

// ErrUserQuotaReached if an user reach its quota
var ErrUserQuotaReached = errors.New("user quota reached")
