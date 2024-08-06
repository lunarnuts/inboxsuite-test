package errors

import (
	"errors"
	"fmt"
)

func New(text string) error {
	return errors.New(text)
}

func Wrap(err error, text string) error {
	return fmt.Errorf("%w: %s", err, text)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
