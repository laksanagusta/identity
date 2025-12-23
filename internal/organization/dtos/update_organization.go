package dtos

import (
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

type UpdateOrganizationReq struct {
	OrganizationUUID string              `params:"organizationUUID"`
	Name             nullable.NullString `json:"name"`
	Address          nullable.NullString `json:"address"`
	Latitude         nullable.NullString `json:"latitude"`
	Longitude        nullable.NullString `json:"longitude"`
}

func (r UpdateOrganizationReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Length(1, 255)),
		validation.Field(&r.Address, validation.Length(1, 255)),
		validation.Field(&r.Latitude, is.Latitude),
		validation.Field(&r.Longitude, is.Longitude),
		validation.Field(&r.OrganizationUUID, validation.Required, is.UUIDv4),
	)
}

func (r UpdateOrganizationReq) NewUpdateOrganization(cred entities.AuthenticatedUser) entities.Organization {
	organization := entities.Organization{
		Name:      r.Name,
		Address:   r.Address,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}

	organization.UUID = r.OrganizationUUID
	organization.UpdateModel(cred.Username)

	return organization
}
