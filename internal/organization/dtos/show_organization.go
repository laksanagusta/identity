package dtos

import (
	"time"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"
)

type ShowOrganizationRes struct {
	UUID    string              `json:"id"`
	Name    nullable.NullString `json:"name"`
	Address nullable.NullString `json:"address"`
	Type    nullable.NullString `json:"type"`

	Organizations []ShowOrganizationRes `json:"organizations"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
}

func NewShowOrganizationRes(organization *entities.Organization) ShowOrganizationRes {
	if organization == nil {
		return ShowOrganizationRes{}
	}

	res := ShowOrganizationRes{
		UUID:      organization.UUID,
		Name:      organization.Name,
		Address:   organization.Address,
		Type:      organization.Type,
		CreatedAt: organization.CreatedAt,
		CreatedBy: organization.CreatedBy,
	}

	if organization.Children != nil {
		for _, child := range organization.Children {
			if child != nil {
				res.Organizations = append(res.Organizations, NewShowOrganizationRes(child))
			}
		}
	}

	return res
}
