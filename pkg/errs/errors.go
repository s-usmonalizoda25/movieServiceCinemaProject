package errs

import "errors"

var (
	ErrInvalidAgeLimit  = errors.New("invalid age limit")
	ErrTitleEmpty       = errors.New("title cannot be empty")
	ErrDurationInvalid  = errors.New("duration must be a positive number")
	ErrDescriptionShort = errors.New("description is too short")
	ErrInternalServer   = errors.New("internal server error")
)
