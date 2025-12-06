package dtos

import (
	"github.com/invopop/validation"
	"github.com/laksanagusta/identity/internal/entities"
)

type CreateUserRoleReq struct {
	UserUUID string `json:"user_id"`
	RoleUUID string `json:"role_id"`
}

func (r CreateUserRoleReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserUUID, validation.Required, validation.Length(1, 36)),
		validation.Field(&r.RoleUUID, validation.Required, validation.Length(1, 36)),
	)
}

func (r CreateUserRoleReq) NewUserRole(username string) entities.UserRole {
	userRole := entities.UserRole{
		UserUUID: r.UserUUID,
		RoleUUID: r.RoleUUID,
	}

	userRole.BaseModel = entities.NewBaseModel(username)

	return userRole
}
