package kensho

import (
	"fmt"
	"reflect"
)

func validConstraint(_ *ValidationContext) error {
	return nil
}

func StructConstraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	t := reflect.TypeOf(ctx.Value())
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		ctx.BuildViolation("not_struct", nil).AddViolation()
	}

	return nil
}

func StringConstraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	t := reflect.TypeOf(ctx.Value())
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	if t.Kind() != reflect.String {
		ctx.BuildViolation("not_string", nil).AddViolation()
	}

	return nil
}

func RequiredConstraint(ctx *ValidationContext) error {
	if ctx.Value() != nil {
		switch reflect.TypeOf(ctx.Value()).Kind() {
		case reflect.String:
			if len(ctx.Value().(string)) > 0 {
				return nil
			}
		case reflect.Array, reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface:
			if !reflect.ValueOf(ctx.Value()).IsNil() {
				return nil
			}
		default:
			return nil
		}
	}

	ctx.BuildViolation("is_required", nil).AddViolation()

	return nil
}

func LengthConstraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	var length int
	switch v := ctx.Arg().(type) {
	case int:
		length = v
	case int64:
		length = int(v)
	case float64:
		length = int(v)
	default:
		return fmt.Errorf("invalid argument to length: expected int, int64, or float64, got %T", v)
	}

	val := reflect.ValueOf(ctx.Value())
	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		if val.Len() != length {
			ctx.BuildViolation("invalid_length", map[string]interface{}{
				"length": length,
			}).AddViolation()
		}
		return nil
	default:
		return fmt.Errorf("expected a slice, map, or string value, got %T", ctx.Value())
	}
}

func MinConstraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	var min int
	switch v := ctx.Arg().(type) {
	case int:
		min = v
	case int64:
		min = int(v)
	case float64:
		min = int(v)
	default:
		// Return a descriptive error instead of panicking
		return fmt.Errorf("invalid argument to min: expected int, int64, or float64, got %T", v)
	}

	val := reflect.ValueOf(ctx.Value())
	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:

		length := val.Len()
		if length < min {
			ctx.BuildViolation("too_short", map[string]interface{}{
				"min":    min,
				"length": length,
			}).AddViolation()
		}

		return nil
	default:
		return fmt.Errorf("expected a slice, map, or string value, got %T", ctx.Value())
	}
}

func MaxConstraint(ctx *ValidationContext) error {
	if ctx.Value() == nil {
		return nil
	}

	var max int
	switch v := ctx.Arg().(type) {
	case int:
		max = v
	case int64:
		max = int(v)
	case float64:
		max = int(v)
	default:
		return fmt.Errorf("invalid argument to max: expected int, int64, or float64, got %T", v)
	}

	val := reflect.ValueOf(ctx.Value())
	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		length := val.Len()
		if length > max {
			ctx.BuildViolation("too_long", map[string]interface{}{
				"max":    max,
				"length": length,
			}).AddViolation()
		}
		return nil
	default:
		return fmt.Errorf("expected a slice, map, or string value, got %T", ctx.Value())
	}
}
