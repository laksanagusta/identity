package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://dikalaksana.com,https://dikalaksana.com,http://localhost:5173",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Authorization,Content-Type,Traceparent,Accept-Encoding",
		AllowCredentials: true,
	})
}
