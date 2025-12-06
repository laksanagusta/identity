package v1

import (
	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/user"

	"github.com/gofiber/fiber/v2"
)

func MapUser(routes fiber.Router, public fiber.Router, h user.Handlers) {
	public.Post("/login", h.Login)
	public.Post("/register", h.Create)

	userGroup := routes.Group("/users")
	userGroup.Delete("/:userUUID", h.Delete)
	userGroup.Get("/", h.Index)
	userGroup.Post("/login", h.Login)
	userGroup.Patch("/:userUUID", h.Update)
	userGroup.Get("/whoami", h.Whoami)
	userGroup.Get("/:userId", h.Show)
	userGroup.Patch("/:userUUID/change-password", h.ChangePassword)

	roleGroup := routes.Group("/roles")
	roleGroup.Get("/", h.Role)
	roleGroup.Post("/", h.CreateRole)
	roleGroup.Delete("/:roleUUID", h.DeleteRole)

	userRoleGroup := routes.Group("/user-roles")
	userRoleGroup.Post("/", h.CreateUserRole)
	userRoleGroup.Delete("/:userRoleUUID", h.DeleteUserRole)

	permissionGroup := routes.Group("/permissions")
	permissionGroup.Post("/", h.CreatePermission)
	permissionGroup.Patch("/:permissionUUID", h.UpdatePermission)
	permissionGroup.Delete(":permissionUUID", h.DeletePermission)
	permissionGroup.Get("/", h.IndexPermission)

	rolePermissionGroup := routes.Group("/role-permissions")
	rolePermissionGroup.Post("/", h.CreateRolePermission)
	rolePermissionGroup.Delete("/:rolePermissionUUID", h.DeleteRolePermission)
}

// MapExternalUser maps external API routes with API Key authentication
func MapExternalUser(routes fiber.Router, h *ExternalUserHandler, cfg config.Config) {
	// External user endpoints (routes parameter already includes /api/v1/external prefix)
	usersGroup := routes.Group("/users")
	usersGroup.Get("/", h.GetUsers)
	usersGroup.Get("/:id", h.GetUser)
}
