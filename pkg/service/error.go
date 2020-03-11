package service

import "errors"

// ErrCategoryNotFound if a category is not found
var ErrCategoryNotFound = errors.New("category not found")

// ErrUserQuotaReached if an user reach its quota
var ErrUserQuotaReached = errors.New("user quota reached")
