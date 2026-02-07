// Package validation provides programmatic struct validation using ozzo-validation.
package validation

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/samber/oops"
)

// Re-export commonly used validators for convenience.
var (
	Required     = validation.Required
	NotNil       = validation.NotNil
	NilOrNotEmpty = validation.NilOrNotEmpty
	Nil          = validation.Nil
	Empty        = validation.Empty
	Skip         = validation.Skip
	In           = validation.In
	NotIn        = validation.NotIn
	Length       = validation.Length
	RuneLength   = validation.RuneLength
	Min          = validation.Min
	Max          = validation.Max
	Match        = validation.Match
	Date         = validation.Date
	Each         = validation.Each
	When         = validation.When
)

// Re-export is validators for common formats.
var (
	IsEmail        = is.Email
	IsEmailFormat  = is.EmailFormat
	IsURL          = is.URL
	IsRequestURL   = is.RequestURL
	IsRequestURI   = is.RequestURI
	IsAlpha        = is.Alpha
	IsDigit        = is.Digit
	IsAlphanumeric = is.Alphanumeric
	IsUUID         = is.UUID
	IsUUIDv4       = is.UUIDv4
	IsInt          = is.Int
	IsFloat        = is.Float
	IsLowerCase    = is.LowerCase
	IsUpperCase    = is.UpperCase
)

// Validatable is implemented by types that can validate themselves.
type Validatable interface {
	Validate(ctx context.Context) error
}

// Validate validates a Validatable and wraps errors with oops.
func Validate(ctx context.Context, v Validatable) error {
	if err := v.Validate(ctx); err != nil {
		return wrapValidationError(ctx, err)
	}
	return nil
}

// ValidateStruct validates struct fields programmatically.
// Usage:
//
//	err := validation.ValidateStruct(ctx, &input,
//	    validation.Field(&input.Email, validation.Required, validation.IsEmail),
//	    validation.Field(&input.Name, validation.Required, validation.Length(2, 100)),
//	)
func ValidateStruct(ctx context.Context, structPtr any, fields ...*validation.FieldRules) error {
	if err := validation.ValidateStruct(structPtr, fields...); err != nil {
		return wrapValidationError(ctx, err)
	}
	return nil
}

// Field creates a validation rule for a struct field.
func Field(fieldPtr any, rules ...validation.Rule) *validation.FieldRules {
	return validation.Field(fieldPtr, rules...)
}

// wrapValidationError wraps validation errors with oops context.
func wrapValidationError(ctx context.Context, err error) error {
	return oops.
		WithContext(ctx).
		Wrap(err)
}
