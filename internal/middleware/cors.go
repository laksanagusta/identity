package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://dikalaksana.com,https://dikalaksana.com,http://localhost:5173,https://api.marvcore.com,https://marvcore.com,https://www.marvcore.com,http://localhost:5002,http://localhost:3001,https://orion.marvcore.com",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Authorization,Content-Type,Accept,Accept-Language,Content-Length,Accept-Encoding,X-Requested-With,Traceparent,X-CSRF-Token,Cache-Control,Pragma,x-api-key",
		AllowCredentials: true,
	})
}
