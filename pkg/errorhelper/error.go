package errorhelper

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrForbiddenAccess     = errors.New("forbidden access")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrConflict            = errors.New("conflict")
	ErrTimeout             = errors.New("timeout")
	ErrBadRequest          = errors.New("bad request")
	ErrGateway             = errors.New("gateway")
	ErrInternalServer      = errors.New("internal server error")
)

type AppError struct {
	Err     error
	Message string
	errMap  any
}

type TaskStockErr struct {
	UUID        string   `json:"id"`
	ProductUUID []string `json:"product_id,omitempty"`
}

func (h AppError) Error() string {
	return h.Err.Error()
}

func ErrMap[T any](h *AppError) map[string][]T {
	return h.errMap.(map[string][]T)
}

func BadRequest(message string) error {
	return &AppError{
		Err:     ErrBadRequest,
		Message: message,
	}
}
func BadRequestMap[T any](errMap map[string][]T) error {
	return &AppError{
		Err:     ErrBadRequest,
		Message: "bad_request",
		errMap:  errMap,
	}
}

func Unauthorized() error {
	return &AppError{
		Message: "unauthorized",
		Err:     ErrUnauthorized,
	}
}

func UnauthorizedWithMessage(message string) error {
	return &AppError{
		Message: message,
		Err:     ErrUnauthorized,
	}
}

func Forbidden() error {
	return &AppError{
		Message: "forbidden",
		Err:     ErrForbiddenAccess,
	}
}

func ForbiddenWithMessage(message string) error {
	return &AppError{
		Message: message,
		Err:     ErrForbiddenAccess,
	}
}

func ForbiddenMap(errMap map[string][]string) error {
	return &AppError{
		errMap:  errMap,
		Message: "forbidden",
		Err:     ErrForbiddenAccess,
	}
}

func NotFound() error {
	return &AppError{
		Message: "not found",
		Err:     ErrNotFound,
	}
}

func NotFoundWithMessage(message string) error {
	return &AppError{
		Message: message,
		Err:     ErrNotFound,
	}
}

func Conflict() error {
	return &AppError{
		Message: "error conflict! must rollback",
		Err:     ErrConflict,
	}
}

func GatewayTimeout() error {
	return &AppError{
		Message: "gateway timeout",
		Err:     ErrGateway,
	}
}

func UnexpectedUnmarshal(statusCode int, body []byte) error {
	return fmt.Errorf("Status Code: %d | Body: %s", statusCode, string(body))
}
