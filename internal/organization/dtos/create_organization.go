package dtos

import (
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type CreateNewOrganizationReq struct {
	Name      nullable.NullString `json:"name"`
	Address   nullable.NullString `json:"address"`
	Latitude  nullable.NullString `json:"latitude"`
	Longitude nullable.NullString `json:"longitude"`
	Type      nullable.NullString `json:"type"`
	ParentId  nullable.NullString `json:"parent_id"`
}

func (r CreateNewOrganizationReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Address, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Latitude, validation.Required, is.Latitude),
		validation.Field(&r.Longitude, validation.Required, is.Longitude),
		validation.Field(&r.Latitude, validation.Required),
		validation.Field(&r.Longitude, validation.Required),
		validation.Field(&r.Type, validation.Required, validation.Length(1, 255)),
	)
}

func (r CreateNewOrganizationReq) NewOrganization(cred entities.AuthenticatedUser) entities.Organization {
	organization := entities.Organization{
		Name:       r.Name,
		Address:    r.Address,
		Type:       r.Type,
		ParentUUID: r.ParentId,
	}

	baseModel := entities.NewBaseModel(cred.Username)
	organization.BaseModel = baseModel

	return organization
}
