package v1

import (
	"net/http"

	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/middleware"
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

	// Safely get authenticated user
	authUser, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.userUc.Update(
		c.Context(),
		*authUser,
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
	// Safely get authenticated user
	authUser, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

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

func (h *userHandler) IndexRole(c *fiber.Ctx) error {
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

	roles, pagination, err := h.userUc.IndexRole(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	res := dtos.NewListRoleResp2(roles)
	pagination.Data = res

	return c.JSON(pagination)
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

	// Safely get authenticated user
	authUser, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.userUc.Delete(
		c.Context(),
		*authUser,
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

	// Safely get authenticated user
	authUser, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.userUc.ChangePassword(
		c.Context(),
		*authUser,
		changePassword,
	)
	if err != nil {
		return err
	}

	return c.SendStatus(http.StatusOK)
}

func (h *userHandler) ApproveUser(c *fiber.Ctx) error {
	var params struct {
		UserUUID string `params:"userUUID"`
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

	err = h.userUc.ApproveUser(
		c.Context(),
		*authUser,
		params.UserUUID,
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(entities.ResponseData{
		Data: map[string]string{"message": "User berhasil disetujui"},
	})
}

func (h *userHandler) RejectUser(c *fiber.Ctx) error {
	var params struct {
		UserUUID string `params:"userUUID"`
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

	err = h.userUc.RejectUser(
		c.Context(),
		*authUser,
		params.UserUUID,
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(entities.ResponseData{
		Data: map[string]string{"message": "User berhasil ditolak"},
	})
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

	// Safely get authenticated user
	cred, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	roleID, err := h.userUc.CreateRole(
		c.Context(),
		createRole,
		*cred,
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(entities.ResponseData{Data: map[string]string{"id": roleID}})
}

func (h *userHandler) ShowRole(c *fiber.Ctx) error {
	var params struct {
		RoleUUID string `params:"roleUUID"`
	}

	err := c.ParamsParser(&params)
	if err != nil {
		return err
	}

	role, err := h.userUc.ShowRole(
		c.Context(),
		params.RoleUUID,
	)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(entities.ResponseData{Data: dtos.NewShowRoleResp(role)})
}

func (h *userHandler) UpdateRole(c *fiber.Ctx) error {
	var updateRole dtos.UpdateRoleReq
	err := c.ParamsParser(&updateRole)
	if err != nil {
		return err
	}

	err = c.BodyParser(&updateRole)
	if err != nil {
		return err
	}

	err = updateRole.Validate()
	if err != nil {
		return err
	}

	// Safely get authenticated user
	cred, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.userUc.UpdateRole(
		c.Context(),
		updateRole,
		*cred,
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

	// Safely get authenticated user
	cred, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

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

	// Safely get authenticated user
	cred, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

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

	// Safely get authenticated user
	cred, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

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

	// Safely get authenticated user
	cred, err := middleware.GetAuthenticatedUser(c)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

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
