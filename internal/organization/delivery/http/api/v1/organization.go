package v1

import (
	"net/http"

	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/middleware"
	"github.com/laksanagusta/identity/internal/organization"
	"github.com/laksanagusta/identity/internal/organization/dtos"

	"github.com/gofiber/fiber/v2"
)

func NewOrganizationHandler(config config.Config, organizationUc organization.UseCase) organization.Handlers {
	return &organizationHandler{
		config:         config,
		organizationUc: organizationUc,
	}
}

type organizationHandler struct {
	config         config.Config
	organizationUc organization.UseCase
}

func (h *organizationHandler) Organization(c *fiber.Ctx) error {
	var createOrganization dtos.CreateNewOrganizationReq
	err := c.BodyParser(&createOrganization)
	if err != nil {
		return err
	}

	err = createOrganization.Validate()
	if err != nil {
		return err
	}

	// Safely get authenticated user
	authUser, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	organizationId, err := h.organizationUc.Create(
		c.Context(),
		*authUser,
		createOrganization,
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(entities.ResponseData{Data: map[string]string{"organization_id": organizationId}})
}

func (h *organizationHandler) Show(c *fiber.Ctx) error {
	var params struct {
		OrganizationUUID string `params:"organizationUUID"`
	}

	err := c.ParamsParser(&params)
	if err != nil {
		return err
	}

	// Safely get authenticated user
	authUser, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	organization, err := h.organizationUc.Show(
		c.Context(),
		*authUser,
		params.OrganizationUUID,
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(entities.ResponseData{Data: dtos.NewShowOrganizationRes(organization)})
}

func (h *organizationHandler) Index(c *fiber.Ctx) error {
	var listProductReq dtos.ListOrganizationReq
	err := c.QueryParser(&listProductReq)
	if err != nil {
		return err
	}

	err = listProductReq.Validate()
	if err != nil {
		return err
	}

	// Safely get authenticated user
	authUser, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	organizations, metadata, err := h.organizationUc.ListOrganization(
		c.Context(),
		*authUser,
		listProductReq,
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(
		dtos.NewListOrganizationResp(organizations, metadata),
	)
}

func (h *organizationHandler) Update(c *fiber.Ctx) error {
	var updateOrganization dtos.UpdateOrganizationReq
	err := c.BodyParser(&updateOrganization)
	if err != nil {
		return err
	}

	err = c.ParamsParser(&updateOrganization)
	if err != nil {
		return err
	}

	err = updateOrganization.Validate()
	if err != nil {
		return err
	}

	// Safely get authenticated user
	authUser, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.organizationUc.Update(
		c.Context(),
		*authUser,
		updateOrganization,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *organizationHandler) Delete(c *fiber.Ctx) error {
	var params struct {
		OrganizationUUID string `params:"organizationUUID"`
	}

	err := c.ParamsParser(&params)
	if err != nil {
		return err
	}

	// Safely get authenticated user
	authUser, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.organizationUc.Delete(
		c.Context(),
		*authUser,
		params.OrganizationUUID,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}
