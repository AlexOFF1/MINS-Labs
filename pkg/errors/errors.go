package errors

import "fmt"

type AppError struct {
	Op   string
	Code string
	Err  error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Code, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Code)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewNotFoundError(op, entity string) *AppError {
	return &AppError{
		Op:   op,
		Code: "NOT_FOUND",
		Err:  fmt.Errorf("%s not found", entity),
	}
}

func NewValidationError(op, field, reason string) *AppError {
	return &AppError{
		Op:   op,
		Code: "VALIDATION_ERROR",
		Err:  fmt.Errorf("invalid field '%s': %s", field, reason),
	}
}

func NewDuplicateError(op, entity, key string) *AppError {
	return &AppError{
		Op:   op,
		Code: "DUPLICATE_ENTRY",
		Err:  fmt.Errorf("%s with key '%s' already exists", entity, key),
	}
}

func NewInternalError(op string, err error) *AppError {
	return &AppError{
		Op:   op,
		Code: "INTERNAL_ERROR",
		Err:  err,
	}
}
