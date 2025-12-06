package organization

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	Organization(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
	Index(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}
