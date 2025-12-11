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
	organizationsGroup.Get("/", h.GetOrganizations)
	organizationsGroup.Get("/:id", h.GetOrganization)
}

// MapPublicOrganization maps public API routes without authentication
func MapPublicOrganization(routes fiber.Router, h *PublicOrganizationHandler, cfg config.Config) {
	// Public organization endpoints (routes parameter already includes /api/public/v1 prefix)
	organizationsGroup := routes.Group("/organizations")
	organizationsGroup.Get("/", h.GetOrganizations)
	organizationsGroup.Get("/:id", h.GetOrganization)
}
