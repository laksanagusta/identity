package dtos

import (
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"
)

type WhoamiRes struct {
	UUID                string                `json:"id"`
	EmployeeID          nullable.NullString   `json:"employee_id"`
	Username            nullable.NullString   `json:"username"`
	FirstName           nullable.NullString   `json:"first_name"`
	LastName            nullable.NullString   `json:"last_name"`
	PhoneNumber         nullable.NullString   `json:"phone_number"`
	AvatarGradientStart nullable.NullString   `json:"avatar_gradient_start"`
	AvatarGradientEnd   nullable.NullString   `json:"avatar_gradient_end"`
	Roles               []WhoamiResRole       `json:"roles"`
	Permissions         []WhoamiResPermission `json:"permissions"`
	Organization        WhoamiResOrganization `json:"organization"`
	Scopes              []string              `json:"scopes"`
}

type WhoamiResRole struct {
	UUID string              `json:"id"`
	Name nullable.NullString `json:"name"`
}

type WhoamiResOrganization struct {
	UUID string              `json:"id"`
	Name nullable.NullString `json:"name"`
	Type nullable.NullString `json:"type"`
}

type WhoamiResPermission struct {
	UUID     string              `json:"id"`
	Name     nullable.NullString `json:"name"`
	Resource nullable.NullString `json:"resource"`
	Action   nullable.NullString `json:"action"`
}

func NewWhoamiRes(user entities.User, scopes []string) WhoamiRes {
	whoami := WhoamiRes{
		UUID:                user.UUID,
		EmployeeID:          user.EmployeeID,
		Username:            user.Username,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		PhoneNumber:         user.PhoneNumber,
		AvatarGradientStart: user.AvatarGradientStart,
		AvatarGradientEnd:   user.AvatarGradientEnd,
		Organization: WhoamiResOrganization{
			UUID: user.Organization.UUID,
			Name: user.Organization.Name,
			Type: user.Organization.Type,
		},
		Scopes: scopes,
	}

	for _, role := range user.Roles {
		whoami.Roles = append(whoami.Roles, WhoamiResRole{
			UUID: role.UUID,
			Name: role.Name,
		})
	}

	for _, permission := range user.Permissions {
		whoami.Permissions = append(whoami.Permissions, WhoamiResPermission{
			UUID:     permission.UUID,
			Name:     permission.Name,
			Action:   permission.Action,
			Resource: permission.Resource,
		})
	}

	return whoami
}
