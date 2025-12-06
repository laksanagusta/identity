package dtos

import (
	"github.com/invopop/validation"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"
)

type UpdatePermissionReq struct {
	UUID     string              `params:"permissionUUID"`
	Name     nullable.NullString `json:"name"`
	Action   nullable.NullString `json:"action"`
	Resource nullable.NullString `json:"resource"`
}

func (r UpdatePermissionReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UUID, validation.Required, validation.Length(1, 36)),
		validation.Field(&r.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Action, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Resource, validation.Required, validation.Length(1, 255)),
	)
}

func (r UpdatePermissionReq) NewPermission(username string) entities.Permission {
	permission := entities.Permission{
		Name:     r.Name,
		Action:   r.Action,
		Resource: r.Resource,
	}

	permission.BaseModel.UUID = r.UUID
	permission.BaseModel.UpdateModel(username)
	return permission
}
