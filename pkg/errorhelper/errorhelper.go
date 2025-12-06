package errorhelper

import (
	// "encoding/json"

	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/laksanagusta/identity/constants"

	"github.com/golang-jwt/jwt/v4"
	"github.com/invopop/validation"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func HttpHandleError(c *fiber.Ctx, err error) error {
	fmt.Println(err)
	// Debug Mode for local env
	// if os.Getenv("APP_ENV") == "local" {
	// 	log.Println(err)
	// }

	// Unauthorized Error
	var jwtErr *jwt.ValidationError
	if errors.As(err, &jwtErr) || err.Error() == "Missing or malformed JWT" {
		return c.Status(http.StatusUnauthorized).JSON(Error{
			Message: constants.UnauthorizedError,
		})
	}

	// Path Parse Error
	var numErr *strconv.NumError
	if errors.As(err, &numErr) {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: constants.MalformedBodyError,
		})
	}

	// Handle Http Error
	var appErr *AppError
	if errors.As(err, &appErr) {
		if errors.Is(appErr.Err, ErrBadRequest) {
			if appErr.errMap != nil {
				return c.Status(http.StatusBadRequest).JSON(Error{
					Message: appErr.Message,
					Errors:  appErr.errMap,
				})
			}

			return c.Status(http.StatusBadRequest).JSON(Error{
				Message: appErr.Message,
			})
		}

		if errors.Is(appErr.Err, ErrUnauthorized) {
			return c.Status(http.StatusUnauthorized).JSON(Error{
				Message: appErr.Message,
			})
		}

		if errors.Is(appErr.Err, ErrForbiddenAccess) {
			if appErr.errMap != nil {
				return c.Status(http.StatusForbidden).JSON(Error{
					Message: appErr.Message,
					Errors:  appErr.errMap,
				})
			}

			return c.Status(http.StatusForbidden).JSON(Error{
				Message: appErr.Message,
			})
		}

		if errors.Is(appErr.Err, ErrNotFound) {
			return c.Status(http.StatusNotFound).JSON(Error{
				Message: appErr.Message,
			})
		}

		if errors.Is(appErr.Err, ErrConflict) {
			return c.Status(http.StatusConflict).JSON(Error{
				Message: appErr.Message,
			})
		}

		if errors.Is(appErr.Err, ErrGateway) {
			return c.Status(http.StatusGatewayTimeout).JSON(Error{
				Message: appErr.Message,
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(Error{
			Message: appErr.Message,
		})
	}
	var validatorError validation.Errors
	if errors.As(err, &validatorError) {
		mapErr := validationErrorMapping(validatorError)
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: constants.ValidationError,
			Errors:  mapErr,
		})
	}

	// JSON Format Error
	var jsonSyntaxErr *json.SyntaxError
	if errors.As(err, &jsonSyntaxErr) {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: constants.MalformedBodyError,
		})
	}

	// Unmarshal Error
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		var translatedType string
		switch unmarshalErr.Type.Name() {
		// REGEX *int*
		case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
			translatedType = "number"
		case "Time":
			translatedType = "date time"
		case "string":
			translatedType = "string"
		}

		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: constants.MalformedBodyError,
			Errors: map[string][]string{
				unmarshalErr.Field: {fmt.Sprintf("the field must be a valid %s", translatedType)},
			},
		})
	}

	// time parse error
	var timeParseErr *time.ParseError
	if errors.As(err, &timeParseErr) {
		v := timeParseErr.Value
		if v == "" {
			v = "empty string (``)"
		}
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: fmt.Sprintf("invalid time format on %s", v),
		})
	}

	// Query Parameter Error
	var fiberMultiErr fiber.MultiError
	if errors.As(err, &fiberMultiErr) {
		validationErrors := make(map[string][]string)

		for key, err := range fiberMultiErr {
			validationErrors[key] = append(validationErrors[key], err.Error())
		}
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: constants.MalformedQueryError,
			Errors:  validationErrors,
		})
	}

	// Multipart Error
	if errors.Is(err, fasthttp.ErrNoMultipartForm) {
		return c.Status(http.StatusBadRequest).JSON(Error{
			Message: "invalid multipart content-type",
		})
	}

	// Default Fiber Error
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return c.Status(fiberErr.Code).JSON(Error{
			Message: "",
		})
	}

	// TCP connection error
	var tcpErr *net.OpError
	if errors.As(err, &tcpErr) {
		log.Panicf("unable to get tcp connection from %s, shutting down...", tcpErr.Addr.String())
	}

	return c.Status(http.StatusInternalServerError).JSON(Error{
		Message: ErrInternalServer.Error(),
	})
}

func validationErrorMapping(validatorError validation.Errors) map[string][]string {
	mapErr := make(map[string][]string)
	for key, err := range validatorError {
		if errs, ok := err.(validation.Errors); ok {
			newMap := validationErrorMapping(errs)
			mapErr = mergeMapWithKey(key, mapErr, newMap)
		} else {
			mapErr[key] = append(mapErr[key], err.Error())
		}
	}
	return mapErr
}

func mergeMapWithKey(key string, maps ...map[string][]string) map[string][]string {
	res := make(map[string][]string)
	for _, m := range maps {
		for k, v := range m {
			mergedKey := key + "." + k
			res[mergedKey] = append(res[mergedKey], v...)
		}
	}
	if len(res) == 0 {
		return nil
	}
	return res
}
