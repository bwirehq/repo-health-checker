package github

import (
	"errors"
	"fmt"
)

type ErrorKind string

const (
	ErrNotFound    ErrorKind = "not_found"
	ErrForbidden   ErrorKind = "forbidden"
	ErrRateLimited ErrorKind = "rate_limited"
	ErrUnavailable ErrorKind = "unavailable"
)

type APIError struct {
	Kind    ErrorKind
	Message string
	Err     error
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Kind)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func wrap(kind ErrorKind, format string, err error, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	return &APIError{Kind: kind, Message: msg, Err: err}
}

func IsKind(err error, kind ErrorKind) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.Kind == kind
}
