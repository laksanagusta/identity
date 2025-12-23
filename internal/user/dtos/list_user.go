package dtos

import (
	"time"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"

	"github.com/invopop/validation"
)

type ListUserReq struct {
	Page       int        `query:"page"`
	Limit      int        `query:"limit"`
	Search     *string    `query:"search"`
	EmployeeID *string    `query:"employee_id"`
	Sort       *string    `query:"sort"`
	StartTime  *time.Time `query:"start_time"`
	EndTime    *time.Time `query:"end_time"`
	RoleUUID   *string    `query:"role_uuid"`
	Role       *string    `query:"role"`
}

type ListUserRespData struct {
	UUID                string                       `json:"id"`
	EmployeeID          nullable.NullString          `json:"employee_id"`
	Username            nullable.NullString          `json:"username"`
	Firstname           nullable.NullString          `json:"first_name"`
	Lastname            nullable.NullString          `json:"last_name"`
	PhoneNumber         nullable.NullString          `json:"phone_number"`
	AvatarGradientStart nullable.NullString          `json:"avatar_gradient_start"`
	AvatarGradientEnd   nullable.NullString          `json:"avatar_gradient_end"`
	Roles               []ListUserRespDataRole       `json:"roles"`
	Organization        ListUserRespDataOrganization `json:"organizations"`
	CreatedAt           nullable.NullString          `json:"created_at"`
	IsApproved          bool                         `json:"is_approved"`
}

type ListUserRespDataRole struct {
	UUID        string              `json:"id"`
	Name        nullable.NullString `json:"name"`
	Description nullable.NullString `json:"description"`
}

type ListUserRespDataOrganization struct {
	UUID string              `json:"id"`
	Name nullable.NullString `json:"name"`
}

type ListUserRespMetadata struct {
	Count       int64 `json:"count"`
	TotalCount  int64 `json:"total_count"`
	CurrentPage int64 `json:"current_page"`
	TotalPage   int64 `json:"total_page"`
}

func (r ListUserReq) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Limit, validation.Required, validation.Min(1)),
		validation.Field(&r.Page, validation.Required, validation.Min(1)),
		validation.Field(&r.Search, validation.Length(1, 255)),
		validation.Field(&r.Sort, validation.Length(1, 255)),
		validation.Field(&r.StartTime, validation.Length(1, 255), validation.Date(time.RFC3339)),
		validation.Field(&r.EndTime, validation.Length(1, 255), validation.Date(time.RFC3339)),
		validation.Field(&r.Sort, validation.Length(1, 255)),
	)

	return err
}

func NewListUserResp(users []*entities.User) []ListUserRespData {
	data := make([]ListUserRespData, len(users))
	for i, user := range users {
		data[i] = ListUserRespData{
			UUID:                user.UUID,
			EmployeeID:          user.EmployeeID,
			Username:            user.Username,
			Firstname:           user.FirstName,
			Lastname:            user.LastName,
			PhoneNumber:         user.PhoneNumber,
			AvatarGradientStart: user.AvatarGradientStart,
			AvatarGradientEnd:   user.AvatarGradientEnd,
			Organization: ListUserRespDataOrganization{
				UUID: user.Organization.UUID,
				Name: user.Organization.Name,
			},
			CreatedAt:  nullable.NewString(user.CreatedAt.Format(time.RFC3339)),
			IsApproved: user.IsApproved,
		}

		for _, role := range user.Roles {
			roles := ListUserRespDataRole{
				UUID:        role.UUID,
				Name:        role.Name,
				Description: role.Description,
			}
			data[i].Roles = append(data[i].Roles, roles)
		}
	}

	return data
}
