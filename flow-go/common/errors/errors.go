package errors

import "errors"

func NewErrChainEnd(message string) error {
	return errors.New(message)
}

func NewChainNotFound(message string) error {
	return errors.New(message)
}

func NewSlotNotFound(message string) error {
	return errors.New(message)
}
