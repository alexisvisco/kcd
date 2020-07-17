package hook

import (
	"context"

	"github.com/expectedsh/errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Validate func(ctx context.Context, input interface{}) error

func DefaultValidateHook(ctx context.Context, input interface{}) error {
	var err validation.Errors

	switch v := input.(type) {
	case validation.ValidatableWithContext:
		err = v.ValidateWithContext(ctx).(validation.Errors)
	case validation.Validatable:
		err = v.Validate().(validation.Errors)
	}

	if len(err) == 0 {
		return nil
	}

	return errors.
		NewWithKind(errors.KindInvalidArgument, "the request has one or multiple invalid fields").
		WithMetadata("kcd.fields", err)

	return nil
}
