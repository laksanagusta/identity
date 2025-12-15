package middleware

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRecoveryMiddleware_CatchesPanic(t *testing.T) {
	// Create a test logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugarLogger := logger.Sugar()

	// Create a test Fiber app
	app := fiber.New()

	// Add recovery middleware
	app.Use(RecoveryMiddleware(sugarLogger))

	// Create a route that panics
	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("test panic - aplikasi tidak boleh berhenti!")
	})

	// Create a normal route
	app.Get("/normal", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Test panic route
	req := httptest.NewRequest("GET", "/panic", nil)
	resp, err := app.Test(req, -1)

	// Assert no error from app.Test (panic was recovered)
	assert.NoError(t, err)

	// Assert status is 500 (internal server error)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	// Read and verify response body
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Internal server error occurred")

	// Test normal route still works after panic
	req2 := httptest.NewRequest("GET", "/normal", nil)
	resp2, err2 := app.Test(req2, -1)

	// Assert normal route works fine
	assert.NoError(t, err2)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)
}
