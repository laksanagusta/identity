package dtos

import (
	"time"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"
)

type ShowOrganizationRes struct {
	UUID    string                `json:"id"`
	Name    nullable.NullString   `json:"name"`
	Code    nullable.NullString   `json:"code"`
	Address nullable.NullString   `json:"address"`
	Type    nullable.NullString   `json:"type"`
	Parent  *ShowOrganizationRes  `json:"parent,omitempty"`

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
		Code:      organization.Code,
		Address:   organization.Address,
		Type:      organization.Type,
		CreatedAt: organization.CreatedAt,
		CreatedBy: organization.CreatedBy,
	}

	// Include parent information if available
	if organization.Parent != nil {
		parentRes := NewShowOrganizationRes(organization.Parent)
		res.Parent = &parentRes
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
