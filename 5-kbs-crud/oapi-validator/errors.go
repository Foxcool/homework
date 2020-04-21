package middleware

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type CompositeErrorType int

const (
	CompositeErrorTypeValidation = iota + 1
)

type CompositeError struct {
	Type   CompositeErrorType
	errors []error
}

func NewCompositeError(t CompositeErrorType) *CompositeError {
	return &CompositeError{
		Type:   t,
		errors: []error{},
	}
}

func (ce *CompositeError) Error() string {
	var stringErrors []string
	for _, err := range ce.errors {
		stringErrors = append(stringErrors, err.Error())
	}

	return strings.Join(stringErrors, "\n")
}

func (ce *CompositeError) Errors() []error {
	return ce.errors
}

func (ce *CompositeError) add(err error) {
	ce.errors = append(ce.errors, err)
}

// Schema error
type SchemaError struct {
	Key    string
	Origin *openapi3.SchemaError
}

func (e *SchemaError) Error() string {
	return e.Origin.Error()
}
