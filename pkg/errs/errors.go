package errs

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidAgeLimit  = errors.New("invalid age limit")
	ErrTitleEmpty       = errors.New("title cannot be empty")
	ErrDurationInvalid  = errors.New("duration must be positive")
	ErrDescriptionShort = errors.New("description is too short")
	ErrInternalServer   = errors.New("internal server error")
)

const (
	MsgFailedCreate = "failed to create movie"
	MsgFailedGet    = "failed to get movie"
	MsgFailedUpdate = "failed to update movie"
	MsgFailedDelete = "failed to delete movie"
)

func MapToGRPC(err error) error {
	switch {
	case errors.Is(err, ErrInvalidAgeLimit),
		errors.Is(err, ErrTitleEmpty),
		errors.Is(err, ErrDurationInvalid),
		errors.Is(err, ErrDescriptionShort):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
