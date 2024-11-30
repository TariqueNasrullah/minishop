package domain

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	ErrorMap map[string][]string
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("Validation Errors: %v", ve.ErrorMap)
}

var (
	BadRequestError     = errors.New("bad request")
	NotFoundError       = errors.New("not found")
	InternalServerError = errors.New("internal server error")
)
