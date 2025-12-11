package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/organization"
	"github.com/laksanagusta/identity/internal/organization/dtos/public"
)

// PublicOrganizationHandler handles public API requests for organizations
type PublicOrganizationHandler struct {
	config         config.Config
	organizationUc organization.UseCase
}

// NewPublicOrganizationHandler creates a new public organization handler
func NewPublicOrganizationHandler(cfg config.Config, organizationUc organization.UseCase) *PublicOrganizationHandler {
	return &PublicOrganizationHandler{
		config:         cfg,
		organizationUc: organizationUc,
	}
}

// GetOrganizations handles GET /api/public/v1/organizations
// Endpoint untuk public API mendapatkan list organizations tanpa authentication
func (h *PublicOrganizationHandler) GetOrganizations(c *fiber.Ctx) error {
	// Parse query parameters
	var listOrganizationReq public.ListOrganizationReq
	err := c.QueryParser(&listOrganizationReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	// Validate request
	err = listOrganizationReq.Validate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Create a dummy authenticated user for public API
	// Public API doesn't require authentication
	authUser := entities.AuthenticatedUser{
		ID:       "public-api",
		Username: "public-api",
	}

	// Get organizations from use case
	organizations, metadata, err := h.organizationUc.ListOrganization(
		c.Context(),
		authUser,
		listOrganizationReq.ToInternalReq(),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch organizations: " + err.Error(),
		})
	}

	// Convert to public response
	return c.Status(http.StatusOK).JSON(
		public.NewListOrganizationResp(organizations, metadata),
	)
}

// GetOrganization handles GET /api/public/v1/organizations/{id}
// Endpoint untuk public API mendapatkan detail organization tanpa authentication
func (h *PublicOrganizationHandler) GetOrganization(c *fiber.Ctx) error {
	// Get organization UUID from parameter
	organizationUUID := c.Params("id")
	if organizationUUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Organization ID is required",
		})
	}

	// Create a dummy authenticated user for public API
	// Public API doesn't require authentication
	authUser := entities.AuthenticatedUser{
		ID:       "public-api",
		Username: "public-api",
	}

	// Get organization from use case
	organization, err := h.organizationUc.Show(c.Context(), authUser, organizationUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch organization: " + err.Error(),
		})
	}

	if organization == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Organization not found",
		})
	}

	// Convert to public response and use standard response format
	organizationData := public.NewPublicOrganizationRes(*organization)

	return c.Status(http.StatusOK).JSON(entities.ResponseData{Data: organizationData})
}
