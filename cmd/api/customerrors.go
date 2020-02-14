package main

import "fmt"

// InternalError gathers error in the server logic
type InternalError struct {
	Message string `json:"message"`
}

// NewInternalError returns an InternalError
func NewInternalError(message string) error {
	return &InternalError{message}
}

func (err InternalError) Error() string {
	return err.Message
}

// InvalidInputError gathers errors due to an invalid field in the json body
type InvalidInputError struct {
	Field   string
	Value   string
	Message string
}

// NewInvalidInputError returns an InvalidInputError
func NewInvalidInputError(field, value, message string) error {
	return &InvalidInputError{field, value, message}
}

func (err InvalidInputError) Error() string {
	return fmt.Sprintf("'%s' field '%s' is invalid: '%s'", err.Field, err.Value, err.Message)
}
