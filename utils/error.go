package utils

import (
	"errors"
)

var (
	ErrUuidType  error = errors.New("UUID type error")
	ErrStringMax error = errors.New("String length out of defined")
	ErrStringMin error = errors.New("String length lite than defined")
)
