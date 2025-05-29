package errors

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrTicketNotFound                  = errors.New("ticket not found")
	ErrTicketAlreadyClosedOrInProgress = errors.New("ticket already closed or in progress")
	ErrTicketWrongType                 = errors.New("ticket wrong type")
	ErrTicketWrongStatus               = errors.New("ticket wrong status")
	ErrOperatorNotAssigned             = errors.New("operator not assigned")
	ErrActiveTicketAlreadyExists       = errors.New("user already has an active ticket")
)
