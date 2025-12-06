package dtos

import (
	"github.com/invopop/validation"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"
)

type CreatePermissionReq struct {
	Name     nullable.NullString `json:"name"`
	Action   nullable.NullString `json:"action"`
	Resource nullable.NullString `json:"resource"`
}

func (r CreatePermissionReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Action, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Resource, validation.Required, validation.Length(1, 255)),
	)
}

func (r CreatePermissionReq) NewPermission(username string) entities.Permission {
	permission := entities.Permission{
		Name:     r.Name,
		Action:   r.Action,
		Resource: r.Resource,
	}

	permission.BaseModel = entities.NewBaseModel(username)

	return permission
}
