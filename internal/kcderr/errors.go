package kcderr

import (
	"fmt"
	"reflect"
)

type InputError struct {
	Message   string
	FieldType reflect.Type
	Field     string
	Extractor string
	Err       error
}

func (b *InputError) WithMessage(msg string) *InputError {
	b.Message = msg
	return b
}

func (b *InputError) WithErr(err error) *InputError {
	b.Err = err
	return b
}

// ErrorDescription implements the builtin error interface for InputError.
func (b InputError) Error() string {
	if b.Field != "" && b.FieldType != nil {
		return fmt.Sprintf(
			"binding error on Field '%s' of type '%s': %s: %s",
			b.Field,
			b.FieldType.Name(),
			b.Message,
			b.Err.Error(),
		)
	}
	return fmt.Sprintf("binding error: %s", b.Message)
}

type OutputError struct {
	Err error
}

// ErrorDescription implements the builtin error interface for OutputError.
func (o OutputError) Error() string {
	return "unable to marshal response into json format"
}
