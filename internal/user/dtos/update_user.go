package dtos

import (
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type UpdateUserReq struct {
	UserUUID    string              `params:"userUUID"`
	EmployeeID  nullable.NullString `json:"employee_id"`
	Username    nullable.NullString `json:"username"`
	FirstName   nullable.NullString `json:"first_name"`
	LastName    nullable.NullString `json:"last_name"`
	PhoneNumber nullable.NullString `json:"phone_number"`
	RoleUUIDs   []string            `json:"role_ids"`
	Password    nullable.NullString `json:"password"`
}

func (r UpdateUserReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.EmployeeID, validation.Length(1, 50)),
		validation.Field(&r.Username, validation.Length(1, 255)),
		validation.Field(&r.FirstName, validation.Length(1, 255)),
		validation.Field(&r.LastName, validation.Length(1, 255)),
		validation.Field(&r.PhoneNumber, is.UTFNumeric, validation.Length(8, 12)),
		validation.Field(&r.Password, validation.Length(6, 255)),
		validation.Field(&r.RoleUUIDs, validation.Each(is.UUIDv4)),
	)
}

func (r UpdateUserReq) NewUser(cred entities.AuthenticatedUser) entities.User {
	user := entities.User{
		EmployeeID:  r.EmployeeID,
		Username:    r.Username,
		FirstName:   r.FirstName,
		LastName:    r.LastName,
		PhoneNumber: r.PhoneNumber,
	}

	user.BaseModel.UUID = r.UserUUID

	for _, roleUUID := range r.RoleUUIDs {
		role := &entities.Role{}
		role.UUID = roleUUID

		user.Roles = append(user.Roles, role)
	}

	return user
}
