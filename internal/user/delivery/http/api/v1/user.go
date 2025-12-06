package v1

import (
	"net/http"

	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/user"
	"github.com/laksanagusta/identity/internal/user/dtos"
	"github.com/laksanagusta/identity/pkg/pagination"

	"github.com/gofiber/fiber/v2"
)

func NewUserHandler(config config.Config, userUc user.UseCase) user.Handlers {
	return &userHandler{
		config: config,
		userUc: userUc,
	}
}

type userHandler struct {
	config config.Config
	userUc user.UseCase
}

func (h *userHandler) Create(c *fiber.Ctx) error {
	var createUser dtos.CreateNewUserReq
	err := c.BodyParser(&createUser)
	if err != nil {
		return err
	}

	err = createUser.Validate()
	if err != nil {
		return err
	}

	newUUID, err := h.userUc.Create(
		c.Context(),
		createUser,
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(
		entities.ResponseData{Data: map[string]any{"id": newUUID}},
	)
}

func (h *userHandler) Update(c *fiber.Ctx) error {
	var updateUser dtos.UpdateUserReq
	err := c.ParamsParser(&updateUser)
	if err != nil {
		return err
	}

	err = c.BodyParser(&updateUser)
	if err != nil {
		return err
	}

	err = updateUser.Validate()
	if err != nil {
		return err
	}

	err = h.userUc.Update(
		c.Context(),
		*c.Locals("authenticatedUser").(*entities.AuthenticatedUser),
		updateUser,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) Show(c *fiber.Ctx) error {
	var param struct {
		UserUUID string `params:"userId"`
	}
	err := c.ParamsParser(&param)
	if err != nil {
		return err
	}

	user, _, err := h.userUc.Show(
		c.Context(),
		param.UserUUID,
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(entities.ResponseData{Data: dtos.NewShowUserRes(*user)})
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	var login dtos.LoginReq
	err := c.BodyParser(&login)
	if err != nil {
		return err
	}

	err = login.Validate()
	if err != nil {
		return err
	}

	token, err := h.userUc.Login(
		c.Context(),
		login,
	)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK).JSON(entities.ResponseData{Data: map[string]string{"token": token}})

	return nil
}

func (h *userHandler) Whoami(c *fiber.Ctx) error {
	authUser := c.Locals("authenticatedUser").(*entities.AuthenticatedUser)

	user, permissionsStr, err := h.userUc.Show(
		c.Context(),
		authUser.ID,
	)
	if err != nil {
		return err
	}

	response := dtos.NewWhoamiRes(*user, permissionsStr)

	return c.Status(http.StatusOK).JSON(entities.ResponseData{Data: response})
}

func (h *userHandler) Role(c *fiber.Ctx) error {
	role, err := h.userUc.Role(
		c.Context(),
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(entities.ResponseData{Data: dtos.NewListRoleResp(role)})
}

func (h *userHandler) Index(c *fiber.Ctx) error {
	queryParams := make(map[string]string)
	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	queryParser := &pagination.QueryParser{}
	params, err := queryParser.Parse(queryParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters: " + err.Error(),
		})
	}

	users, pagination, err := h.userUc.Index(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	res := dtos.NewListUserResp(users)
	pagination.Data = res

	return c.JSON(pagination)
}

func (h *userHandler) Delete(c *fiber.Ctx) error {
	var params struct {
		UserUUID string `params:"userUUID"`
	}

	err := c.ParamsParser(&params)
	if err != nil {
		return err
	}

	err = h.userUc.Delete(
		c.Context(),
		*c.Locals("authenticatedUser").(*entities.AuthenticatedUser),
		params.UserUUID,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) ChangePassword(c *fiber.Ctx) error {
	var changePassword dtos.ChangePassword
	err := c.ParamsParser(&changePassword)
	if err != nil {
		return err
	}

	err = c.BodyParser(&changePassword)
	if err != nil {
		return err
	}
	err = changePassword.Validate()
	if err != nil {
		return err
	}

	err = h.userUc.ChangePassword(
		c.Context(),
		*c.Locals("authenticatedUser").(*entities.AuthenticatedUser),
		changePassword,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) CreateRole(c *fiber.Ctx) error {
	var createRole dtos.CreateRoleReq
	err := c.BodyParser(&createRole)
	if err != nil {
		return err
	}

	err = createRole.Validate()
	if err != nil {
		return err
	}

	cred := c.Locals("authenticatedUser").(*entities.AuthenticatedUser)
	role := createRole.NewRole(*cred)
	err = h.userUc.CreateRole(
		c.Context(),
		role,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) DeleteRole(c *fiber.Ctx) error {
	var params struct {
		RoleUUID string `params:"roleUUID"`
	}

	err := c.ParamsParser(&params)
	if err != nil {
		return err
	}

	err = h.userUc.DeleteRole(
		c.Context(),
		params.RoleUUID,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) CreateUserRole(c *fiber.Ctx) error {
	var createUserRole dtos.CreateUserRoleReq
	err := c.BodyParser(&createUserRole)
	if err != nil {
		return err
	}

	err = createUserRole.Validate()
	if err != nil {
		return err
	}

	cred := c.Locals("authenticatedUser").(*entities.AuthenticatedUser)
	role := createUserRole.NewUserRole(cred.Username)
	err = h.userUc.CreateUserRole(
		c.Context(),
		role,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) DeleteUserRole(c *fiber.Ctx) error {
	var params struct {
		UserRoleUUID string `params:"userRoleUUID"`
	}

	err := c.ParamsParser(&params)
	if err != nil {
		return err
	}

	err = h.userUc.DeleteUserRole(
		c.Context(),
		params.UserRoleUUID,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) CreatePermission(c *fiber.Ctx) error {
	var createPermissionReq dtos.CreatePermissionReq
	err := c.BodyParser(&createPermissionReq)
	if err != nil {
		return err
	}

	err = createPermissionReq.Validate()
	if err != nil {
		return err
	}

	cred := c.Locals("authenticatedUser").(*entities.AuthenticatedUser)
	permission := createPermissionReq.NewPermission(cred.Username)
	err = h.userUc.CreatePermission(
		c.Context(),
		permission,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) UpdatePermission(c *fiber.Ctx) error {
	var params struct {
		PermissionUUID string `params:"permissionUUID"`
	}
	err := c.ParamsParser(&params)
	if err != nil {
		return err
	}

	var updatePermissionReq dtos.UpdatePermissionReq
	err = c.BodyParser(&updatePermissionReq)
	if err != nil {
		return err
	}

	updatePermissionReq.UUID = params.PermissionUUID

	err = updatePermissionReq.Validate()
	if err != nil {
		return err
	}

	cred := c.Locals("authenticatedUser").(*entities.AuthenticatedUser)
	permission := updatePermissionReq.NewPermission(cred.Username)
	err = h.userUc.UpdatePermission(
		c.Context(),
		permission,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) DeletePermission(c *fiber.Ctx) error {
	var params struct {
		PermissionUUID string `params:"permissionUUID"`
	}

	err := c.ParamsParser(&params)
	if err != nil {
		return err
	}

	err = h.userUc.DeletePermission(
		c.Context(),
		params.PermissionUUID,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) CreateRolePermission(c *fiber.Ctx) error {
	var createRolePermissionReq dtos.CreateRolePermissionReq
	err := c.BodyParser(&createRolePermissionReq)
	if err != nil {
		return err
	}

	err = createRolePermissionReq.Validate()
	if err != nil {
		return err
	}

	cred := c.Locals("authenticatedUser").(*entities.AuthenticatedUser)
	rolePermission := createRolePermissionReq.NewRolePermission(cred.Username)
	err = h.userUc.CreateRolePermission(
		c.Context(),
		rolePermission,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) DeleteRolePermission(c *fiber.Ctx) error {
	var params struct {
		RolePermissionUUID string `params:"rolePermissionUUID"`
	}

	err := c.ParamsParser(&params)
	if err != nil {
		return err
	}

	err = h.userUc.DeleteRolePermission(
		c.Context(),
		params.RolePermissionUUID,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) IndexPermission(c *fiber.Ctx) error {
	queryParams := make(map[string]string)
	c.Context().QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	queryParser := &pagination.QueryParser{}
	params, err := queryParser.Parse(queryParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters: " + err.Error(),
		})
	}

	permissions, pagination, err := h.userUc.IndexPermission(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	res := dtos.NewListPermissionResp(permissions)
	pagination.Data = res

	return c.JSON(pagination)
}
