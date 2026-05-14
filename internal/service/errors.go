package service

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("conflict")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrValidation        = errors.New("validation failed")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInactive          = errors.New("resource inactive")
	ErrBadState          = errors.New("invalid state for operation")
	ErrPayment           = errors.New("payment verification failed")
)
