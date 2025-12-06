package dtos

import (
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	"github.com/laksanagusta/identity/internal/entities"
)

// CreateRolePermissionReq is the request DTO for creating a role_permission
type CreateRolePermissionReq struct {
	RoleUUID       string `json:"role_id"`
	PermissionUUID string `json:"permission_id"`
}

func (r CreateRolePermissionReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RoleUUID, validation.Required, validation.Length(1, 36), is.UUIDv4),
		validation.Field(&r.PermissionUUID, validation.Required, validation.Length(1, 36), is.UUIDv4),
	)
}

func (r CreateRolePermissionReq) NewRolePermission(username string) entities.RolaPermission {
	rolePermission := entities.RolaPermission{
		RoleUUID:       r.RoleUUID,
		PermissionUUID: r.PermissionUUID,
	}
	rolePermission.BaseModel = entities.NewBaseModel(username)
	return rolePermission
}
