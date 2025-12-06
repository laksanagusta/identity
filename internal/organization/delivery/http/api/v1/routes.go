package v1

import (
	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/organization"

	"github.com/gofiber/fiber/v2"
)

func MapOrganization(routes fiber.Router, h organization.Handlers) {
	organizationGroup := routes.Group("/organizations")
	organizationGroup.Post("/", h.Organization)
	organizationGroup.Get("/:organizationUUID", h.Show)
	organizationGroup.Get("/", h.Index)
	organizationGroup.Patch("/:organizationUUID", h.Update)
	organizationGroup.Delete("/:organizationUUID", h.Delete)
}

// MapExternalOrganization maps external API routes with API Key authentication
func MapExternalOrganization(routes fiber.Router, h *ExternalOrganizationHandler, cfg config.Config) {
	// External organization endpoints (routes parameter already includes /api/v1/external prefix)
	organizationsGroup := routes.Group("/organizations")
	organizationsGroup.Get("/:id", h.GetOrganization)
}
