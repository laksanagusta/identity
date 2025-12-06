package v1

import (
	"net/http"

	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/organization"
	"github.com/laksanagusta/identity/internal/organization/dtos/external"

	"github.com/gofiber/fiber/v2"
)

// ExternalOrganizationHandler handles external API requests for organizations
type ExternalOrganizationHandler struct {
	config         config.Config
	organizationUc organization.UseCase
}

// NewExternalOrganizationHandler creates a new external organization handler
func NewExternalOrganizationHandler(cfg config.Config, organizationUc organization.UseCase) *ExternalOrganizationHandler {
	return &ExternalOrganizationHandler{
		config:         cfg,
		organizationUc: organizationUc,
	}
}

// GetOrganization handles GET /api/v1/external/organizations/{id}
// Endpoint untuk external API mendapatkan detail organization dengan API Key authentication
func (h *ExternalOrganizationHandler) GetOrganization(c *fiber.Ctx) error {
	// Get organization UUID from parameter
	organizationUUID := c.Params("id")
	if organizationUUID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Organization ID is required",
		})
	}

	// Create a dummy authenticated user for external API
	// External API uses API Key authentication, not JWT
	authUser := entities.AuthenticatedUser{
		ID:       "external-api",
		Username: "external-api",
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

	// Convert to external response and use standard response format
	organizationData := external.NewExternalOrganizationRes(*organization)

	return c.Status(http.StatusOK).JSON(entities.ResponseData{Data: organizationData})
}
