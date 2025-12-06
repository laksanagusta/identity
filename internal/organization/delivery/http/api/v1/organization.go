package v1

import (
	"net/http"

	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/entities"
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

	organizationId, err := h.organizationUc.Create(
		c.Context(),
		*c.Locals("authenticatedUser").(*entities.AuthenticatedUser),
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

	organization, err := h.organizationUc.Show(
		c.Context(),
		*c.Locals("authenticatedUser").(*entities.AuthenticatedUser),
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

	organizations, metadata, err := h.organizationUc.ListOrganization(
		c.Context(),
		*c.Locals("authenticatedUser").(*entities.AuthenticatedUser),
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

	err = h.organizationUc.Update(
		c.Context(),
		*c.Locals("authenticatedUser").(*entities.AuthenticatedUser),
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

	err = h.organizationUc.Delete(
		c.Context(),
		*c.Locals("authenticatedUser").(*entities.AuthenticatedUser),
		params.OrganizationUUID,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}
