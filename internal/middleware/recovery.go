package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RecoveryMiddleware menangkap panic dan mencegah aplikasi berhenti
func RecoveryMiddleware(logger *zap.SugaredLogger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Ambil stack trace
				stackTrace := string(debug.Stack())

				// Log error dengan detail lengkap
				logger.Errorw("Panic recovered",
					"error", r,
					"path", c.Path(),
					"method", c.Method(),
					"ip", c.IP(),
					"user_agent", c.Get("User-Agent"),
					"stack_trace", stackTrace,
				)

				// Return error response ke client
				err := c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "Internal server error occurred",
					"error":   fmt.Sprintf("%v", r),
				})
				if err != nil {
					logger.Errorw("Failed to send error response", "error", err)
				}
			}
		}()

		return c.Next()
	}
}
