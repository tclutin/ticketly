package errors

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrTicketNotFound      = errors.New("ticket not found")
	ErrTicketAlreadyClosed = errors.New("ticket already exists")
)
