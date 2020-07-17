package kcderr

import (
	"fmt"
	"reflect"
)

type Input struct {
	Message   string
	Type      reflect.Type
	Field     string
	Extractor string
	Err       error
}

func (b *Input) WithMessage(msg string) *Input {
	b.Message = msg
	return b
}

func (b *Input) WithErr(err error) *Input {
	b.Err = err
	return b
}

// ErrorDescription implements the builtin error interface for Input.
func (b Input) Error() string {
	if b.Field != "" && b.Type != nil {
		return fmt.Sprintf(
			"binding error on field '%s' of type '%s': %s",
			b.Field,
			b.Type.Name(),
			b.Message,
		)
	}
	return fmt.Sprintf("binding error: %s", b.Message)
}

type Output struct {
	Err error
}

// ErrorDescription implements the builtin error interface for Output.
func (o Output) Error() string {
	return "unable to marshal response into json format"
}
