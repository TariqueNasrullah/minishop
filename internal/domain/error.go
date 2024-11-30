package domain

import "fmt"

type ValidationError struct {
	ErrorMap map[string][]string
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("Validation Errors: %v", ve.ErrorMap)
}
