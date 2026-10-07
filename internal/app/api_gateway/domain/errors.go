package domain

import "errors"

var (
	ErrNothingInTable        = errors.New("no rows in result set")
	ErrInvalidDataInTable    = errors.New("unable to scan: invalid data in table")
	ErrDuplicatedData        = errors.New("duplicated data in table")
	ErrConstraint            = errors.New("constraint error")
	ErrTooManyRows           = errors.New("too many rows in result set")
	ErrIncorrectRowsAffected = errors.New("unexpected value of 'rows affected'")
	ErrResultNotOk           = errors.New("result not ok")
)
