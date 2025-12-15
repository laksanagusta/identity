package dtos

import (
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"

	"github.com/invopop/validation"
)

type CreateRoleReq struct {
	Name          nullable.NullString `json:"name"`
	Description   nullable.NullString `json:"description"`
	PermissionIDs []string            `json:"permission_ids"`
}

func (r CreateRoleReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Description, validation.Length(0, 255)),
	)
}

func (r CreateRoleReq) NewRole(cred entities.AuthenticatedUser) entities.Role {
	role := entities.Role{
		Name:        r.Name,
		Description: r.Description,
	}

	role.BaseModel = entities.NewBaseModel(cred.Username)

	return role
}

type UpdateRoleReq struct {
	RoleUUID      string              `params:"roleUUID"`
	Name          nullable.NullString `json:"name"`
	Description   nullable.NullString `json:"description"`
	PermissionIDs []string            `json:"permission_ids"`
}

func (r UpdateRoleReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RoleUUID, validation.Required),
	)
}

func (r UpdateRoleReq) NewRole(cred entities.AuthenticatedUser) entities.Role {
	role := entities.Role{
		Name:        r.Name,
		Description: r.Description,
	}

	role.UUID = r.RoleUUID
	role.UpdatedBy = cred.Username

	return role
}
