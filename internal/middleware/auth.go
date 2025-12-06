package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/user"
	"github.com/laksanagusta/identity/pkg/authservice/jwt"
)

// AuthMiddleware creates a middleware that validates JWT tokens locally
func AuthMiddleware(cfg config.Config, userRepo user.Repository) fiber.Handler {
	jwtAuth := jwt.NewJwtAuth(cfg)

	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		// Check Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format. Expected: Bearer <token>",
			})
		}

		token := tokenParts[1]

		// Validate JWT token
		claims, err := jwtAuth.ValidateAndClaimToken(token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid token: %v", err),
			})
		}

		// Extract user ID from claims
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token: missing user_id",
			})
		}

		// Get authenticated user data
		authenticatedUser, err := getUserData(c.Context(), userRepo, userID)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("Authentication failed: %v", err),
			})
		}

		// Store authenticated user in context locals
		c.Locals("authenticatedUser", authenticatedUser)

		// Continue to next handler
		return c.Next()
	}
}

// getUserData retrieves full user data from repository
func getUserData(ctx context.Context, userRepo user.Repository, userUUID string) (*entities.AuthenticatedUser, error) {
	// Find user by UUID
	userEntity, err := userRepo.FindByUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if userEntity == nil {
		return nil, errors.New("user not found")
	}

	// Get user roles
	roles, err := userRepo.FindRoleByUserUUID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Convert roles to AuthRole format
	authRoles := make([]entities.AuthRole, len(roles))
	for i, role := range roles {
		authRoles[i] = entities.AuthRole{
			ID:   role.UUID,
			Name: role.Name.GetOrDefault(),
		}
	}

	// Parse organization UUID
	var orgUUID uuid.UUID
	if userEntity.OrganizationUUID.IsExists {
		orgUUID, err = uuid.Parse(userEntity.OrganizationUUID.GetOrDefault())
		if err != nil {
			return nil, fmt.Errorf("invalid organization UUID format: %w", err)
		}
	}

	// Convert to AuthenticatedUser entity
	authenticatedUser := &entities.AuthenticatedUser{
		ID:          userEntity.UUID,
		EmployeeID:  userEntity.EmployeeID.GetOrDefault(),
		Username:    userEntity.Username.GetOrDefault(),
		FirstName:   userEntity.FirstName.GetOrDefault(),
		LastName:    userEntity.LastName.GetOrDefault(),
		PhoneNumber: userEntity.PhoneNumber.GetOrDefault(),
		Roles:       authRoles,
		Organization: entities.UserOrganization{
			ID:   orgUUID,
			Name: "", // TODO: Join with organization table to get name
			Type: "", // TODO: Join with organization table to get type
		},
	}

	return authenticatedUser, nil
}

// GetAuthenticatedUser retrieves the authenticated user from context
func GetAuthenticatedUser(c *fiber.Ctx) (*entities.AuthenticatedUser, error) {
	user, ok := c.Locals("authenticatedUser").(*entities.AuthenticatedUser)
	if !ok {
		return nil, errors.New("authenticated user not found in context")
	}
	return user, nil
}
