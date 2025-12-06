package v1

import (
	"net/http"

	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/user"
	"github.com/laksanagusta/identity/internal/user/dtos/external"
	"github.com/laksanagusta/identity/pkg/pagination"

	"github.com/gofiber/fiber/v2"
)

// ExternalUserHandler handles external API requests
type ExternalUserHandler struct {
	config config.Config
	userUc user.UseCase
}

// NewExternalUserHandler creates a new external user handler
func NewExternalUserHandler(cfg config.Config, userUc user.UseCase) *ExternalUserHandler {
	return &ExternalUserHandler{
		config: cfg,
		userUc: userUc,
	}
}

// GetUsers handles GET /api/v1/external/users
// Endpoint untuk external API mendapatkan list users dengan API Key authentication
func (h *ExternalUserHandler) GetUsers(c *fiber.Ctx) error {
	// Parse query parameters
	queryParams := make(map[string]string)
	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	// Parse request
	var req external.ExternalListUserReq
	err := c.QueryParser(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters: " + err.Error(),
		})
	}

	// Validate request
	err = req.Validate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Validation error: " + err.Error(),
		})
	}

	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	// Build pagination params
	queryParser := &pagination.QueryParser{}
	params, err := queryParser.Parse(queryParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters: " + err.Error(),
		})
	}

	// Get users from use case
	users, paginationResp, err := h.userUc.Index(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users: " + err.Error(),
		})
	}

	// Convert to external response
	userData := external.NewExternalListUserResp(users)

	// Use same response format as internal API
	paginationResp.Data = userData
	return c.JSON(paginationResp)
}

// GetUser handles GET /api/v1/external/users/{id}
// Endpoint untuk external API mendapatkan detail user dengan API Key authentication
func (h *ExternalUserHandler) GetUser(c *fiber.Ctx) error {
	// Get user UUID from parameter
	userUUID := c.Params("id")
	if userUUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Get user from use case
	user, _, err := h.userUc.Show(c.Context(), userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user: " + err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Convert to external response and use standard response format
	userData := external.NewExternalUserRes(*user)

	return c.Status(http.StatusOK).JSON(entities.ResponseData{Data: userData})
}
