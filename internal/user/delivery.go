package user

import (
	"github.com/gofiber/fiber/v2"
)

type Handlers interface {
	// user
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Whoami(c *fiber.Ctx) error
	Index(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error

	// role
	Role(c *fiber.Ctx) error
	CreateRole(c *fiber.Ctx) error
	DeleteRole(c *fiber.Ctx) error

	// user-role
	CreateUserRole(c *fiber.Ctx) error
	DeleteUserRole(c *fiber.Ctx) error

	// permission
	UpdatePermission(c *fiber.Ctx) error
	DeletePermission(c *fiber.Ctx) error
	CreatePermission(c *fiber.Ctx) error
	IndexPermission(c *fiber.Ctx) error

	// role-permissions
	CreateRolePermission(c *fiber.Ctx) error
	DeleteRolePermission(c *fiber.Ctx) error
}
