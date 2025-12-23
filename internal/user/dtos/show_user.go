package dtos

import (
	"time"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"
)

type ShowUserRes struct {
	Id                  string                  `json:"id"`
	EmployeeID          nullable.NullString     `json:"employee_id"`
	Username            nullable.NullString     `json:"username"`
	FirstName           nullable.NullString     `json:"first_name"`
	LastName            nullable.NullString     `json:"last_name"`
	PhoneNumber         nullable.NullString     `json:"phone_number"`
	AvatarGradientStart nullable.NullString     `json:"avatar_gradient_start"`
	AvatarGradientEnd   nullable.NullString     `json:"avatar_gradient_end"`
	Organization        ShowUserResOrganization `json:"organization"`
	Roles               []ShowUserResRole       `json:"role"`
	CreatedAt           time.Time               `json:"created_at"`
}

type ShowUserResRole struct {
	UUID string              `json:"uuid"`
	Name nullable.NullString `json:"name"`
}

type ShowUserResOrganization struct {
	UUID string              `json:"uuid"`
	Name nullable.NullString `json:"name"`
}

func NewShowUserRes(user entities.User) ShowUserRes {
	res := ShowUserRes{
		Id:                  user.UUID,
		EmployeeID:          user.EmployeeID,
		Username:            user.Username,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		PhoneNumber:         user.PhoneNumber,
		AvatarGradientStart: user.AvatarGradientStart,
		AvatarGradientEnd:   user.AvatarGradientEnd,
		Organization:        ShowUserResOrganization{UUID: user.Organization.UUID, Name: user.Organization.Name},
		CreatedAt:           user.CreatedAt,
	}

	for _, role := range user.Roles {
		res.Roles = append(res.Roles, ShowUserResRole{
			UUID: role.UUID,
			Name: role.Name,
		})
	}

	return res
}
