package external

import (
	"time"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"

	"github.com/invopop/validation"
)

// ExternalListUserReq request untuk external API get list users
type ExternalListUserReq struct {
	Page       int        `query:"page"`
	Limit      int        `query:"limit"`
	Search     *string    `query:"search"`
	EmployeeID *string    `query:"employee_id"`
	Username   *string    `query:"username"`
	IsActive   *bool      `query:"is_active"`
	StartTime  *time.Time `query:"start_time"`
	EndTime    *time.Time `query:"end_time"`
}

func (r ExternalListUserReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Limit, validation.Min(1), validation.Max(100)),
		validation.Field(&r.Page, validation.Min(1)),
		validation.Field(&r.Search, validation.Length(1, 255)),
		validation.Field(&r.Username, validation.Length(1, 255)),
	)
}

// ExternalUserRes response structure untuk external API
type ExternalUserRes struct {
	UUID                string                   `json:"id"`
	EmployeeID          nullable.NullString      `json:"employee_id"`
	Username            nullable.NullString      `json:"username"`
	FirstName           nullable.NullString      `json:"first_name"`
	LastName            nullable.NullString      `json:"last_name"`
	Email               nullable.NullString      `json:"email"`
	PhoneNumber         nullable.NullString      `json:"phone_number"`
	AvatarGradientStart nullable.NullString      `json:"avatar_gradient_start"`
	AvatarGradientEnd   nullable.NullString      `json:"avatar_gradient_end"`
	IsActive            bool                     `json:"is_active"`
	IsApproved          bool                     `json:"is_approved"`
	LastLoginAt         *time.Time               `json:"last_login_at,omitempty"`
	Organization        *ExternalOrganizationRes `json:"organization,omitempty"`
	Roles               []ExternalRoleRes        `json:"roles"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           time.Time                `json:"updated_at"`
}

// ExternalOrganizationRes organization info untuk external API
type ExternalOrganizationRes struct {
	UUID string              `json:"id"`
	Name nullable.NullString `json:"name"`
	Type nullable.NullString `json:"type"`
}

// ExternalRoleRes role info untuk external API
type ExternalRoleRes struct {
	UUID        string              `json:"id"`
	Name        nullable.NullString `json:"name"`
	Description nullable.NullString `json:"description"`
}

// NewExternalUserRes convert dari entities.User ke ExternalUserRes
func NewExternalUserRes(user entities.User) ExternalUserRes {
	res := ExternalUserRes{
		UUID:                user.UUID,
		EmployeeID:          user.EmployeeID,
		Username:            user.Username,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Email:               user.Email,
		PhoneNumber:         user.PhoneNumber,
		AvatarGradientStart: user.AvatarGradientStart,
		AvatarGradientEnd:   user.AvatarGradientEnd,
		IsActive:            user.IsActive,
		IsApproved:          user.IsApproved,
		LastLoginAt:         user.LastLoginAt,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}

	// Add organization info jika ada
	if user.Organization != nil {
		res.Organization = &ExternalOrganizationRes{
			UUID: user.Organization.UUID,
			Name: user.Organization.Name,
			Type: user.Organization.Type,
		}
	}

	// Add roles info jika ada
	if len(user.Roles) > 0 {
		res.Roles = make([]ExternalRoleRes, len(user.Roles))
		for i, role := range user.Roles {
			res.Roles[i] = ExternalRoleRes{
				UUID:        role.UUID,
				Name:        role.Name,
				Description: role.Description,
			}
		}
	}

	return res
}

// NewExternalListUserResp convert dari []*entities.User ke []ExternalUserRes
func NewExternalListUserResp(users []*entities.User) []ExternalUserRes {
	data := make([]ExternalUserRes, len(users))
	for i, user := range users {
		data[i] = NewExternalUserRes(*user)
	}
	return data
}
