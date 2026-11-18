package domain

import "errors"

var (
	ErrNotFound         = errors.New("resource not found")
	ErrConflict         = errors.New("resource conflicts with existing or dependent data")
	ErrInvalidReference = errors.New("referenced resource does not exist")
)

type ValidationError struct{ Message string }

func (e *ValidationError) Error() string { return e.Message }
