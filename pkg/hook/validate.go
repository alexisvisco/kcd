package hook

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate is the default validation hook.
// It use 'ozzo-validation' to validate structure.
// A structure must implement 'ValidatableWithContext' or 'Validatable'.
func Validate(ctx context.Context, input interface{}) error {
	switch v := input.(type) {
	case validation.ValidatableWithContext:
		return v.ValidateWithContext(ctx)
	case validation.Validatable:
		return v.Validate()
	}

	return nil
}
