package server

import (
	"github.com/laksanagusta/identity/internal/middleware"
	organizationhandler "github.com/laksanagusta/identity/internal/organization/delivery/http/api/v1"
	organizationrepository "github.com/laksanagusta/identity/internal/organization/repository"
	organizationusecase "github.com/laksanagusta/identity/internal/organization/usecase"
	userhandler "github.com/laksanagusta/identity/internal/user/delivery/http/api/v1"
	userrepository "github.com/laksanagusta/identity/internal/user/repository"
	userusecase "github.com/laksanagusta/identity/internal/user/usecase"

	"github.com/laksanagusta/identity/pkg/authservice/jwt"
	"github.com/laksanagusta/identity/pkg/database"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) MapHandlers() error {
	check := s.Fiber.Group("/check")
	check.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"version": "2.0",
		})
	})

	apiV1 := s.Fiber.Group("/api/v1")

	apiExternalV1 := s.Fiber.Group("/api/v1/external")

	apiPublicV1 := s.Fiber.Group("/api/public/v1")

	userRepo := userrepository.NewUserRepo(s.DB)
	organizationRepo := organizationrepository.NewOrganizationRepo(s.DB)
	authService := jwt.NewJwtAuth(s.Config)

	txManager := database.NewManager(s.DB)

	// Apply AuthMiddleware to apiV1 group with required dependencies
	apiV1.Use(middleware.AuthMiddleware(s.Config, userRepo))

	userUseCase := userusecase.NewUserUseCase(userusecase.UseCaseParameter{
		UserRepo:         userRepo,
		JwtAuth:          authService,
		OrganizationRepo: organizationRepo,
	})

	organizationUseCase := organizationusecase.NewOrganizationUseCase(organizationusecase.UseCaseParameter{
		OrganizationRepo: organizationRepo,
		TxManager:        txManager,
		UserUC:           userUseCase,
	})
	userHandler := userhandler.NewUserHandler(s.Config, userUseCase)
	userhandler.MapUser(apiV1, apiPublicV1, userHandler)

	// External API routes with API Key authentication
	externalUserHandler := userhandler.NewExternalUserHandler(s.Config, userUseCase)
	userhandler.MapExternalUser(apiExternalV1, externalUserHandler, s.Config)

	organizationHandler := organizationhandler.NewOrganizationHandler(s.Config, organizationUseCase)
	organizationhandler.MapOrganization(apiV1, organizationHandler)

	// External API routes with API Key authentication for organizations
	externalOrganizationHandler := organizationhandler.NewExternalOrganizationHandler(s.Config, organizationUseCase)
	organizationhandler.MapExternalOrganization(apiExternalV1, externalOrganizationHandler, s.Config)

	return nil
}
