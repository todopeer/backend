package resolver

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
)