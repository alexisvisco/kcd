package kcd

import (
	"fmt"
	"reflect"
)

type inputError struct {
	message   string
	fieldType reflect.Type
	field     string
	extractor string
	err       error
}

func (b *inputError) withMessage(msg string) *inputError {
	b.message = msg
	return b
}

func (b *inputError) withErr(err error) *inputError {
	b.err = err
	return b
}

// ErrorDescription implements the builtin error interface for inputError.
func (b inputError) Error() string {
	if b.field != "" && b.fieldType != nil {
		return fmt.Sprintf(
			"binding error on field '%s' of type '%s': %s: %s",
			b.field,
			b.fieldType.Name(),
			b.message,
			b.err.Error(),
		)
	}
	return fmt.Sprintf("binding error: %s", b.message)
}

type outputError struct {
	Err error
}

// ErrorDescription implements the builtin error interface for outputError.
func (o outputError) Error() string {
	return "unable to marshal response into json format"
}
