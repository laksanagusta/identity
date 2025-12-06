package jwt

import (
	"errors"
	"log"
	"time"

	"github.com/laksanagusta/identity/config"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/errorhelper"

	"github.com/golang-jwt/jwt/v5"
)

func NewJwtAuth(config config.Config) JwtAuth {
	return JwtAuth{
		config: config,
	}
}

type JwtAuth struct {
	config config.Config
}

func (s *JwtAuth) GenerateToken(user entities.User, store entities.Organization) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = user.UUID
	claim["username"] = user.Username
	claim["roles"] = user.Roles
	claim["organization_id"] = user.OrganizationUUID
	claim["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(s.config.JWT.SecretKey))
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

func (s *JwtAuth) ValidateAndClaimToken(encodedToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errorhelper.Unauthorized()
		}

		return []byte(s.config.JWT.SecretKey), nil
	})
	if err != nil {
		log.Println("ok", err)
		return nil, err
	}

	if !token.Valid {
		return nil, errorhelper.Unauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token claim error")
	}

	return claims, nil
}
