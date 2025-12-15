package dtos

import (
	"time"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"
)

type ShowOrganizationRes struct {
	UUID    string                 `json:"id"`
	Name    nullable.NullString    `json:"name"`
	Code    nullable.NullString    `json:"code"`
	Address nullable.NullString    `json:"address"`
	Type    nullable.NullString    `json:"type"`
	Parent  *ParentOrganizationRes `json:"parent,omitempty"`

	Organizations []ShowOrganizationRes `json:"organizations"`

	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
}

// ParentOrganizationRes - simplified parent without children to avoid circular reference
type ParentOrganizationRes struct {
	UUID    string              `json:"id"`
	Name    nullable.NullString `json:"name"`
	Code    nullable.NullString `json:"code"`
	Address nullable.NullString `json:"address"`
	Type    nullable.NullString `json:"type"`
}

func NewShowOrganizationRes(organization *entities.Organization) ShowOrganizationRes {
	if organization == nil {
		return ShowOrganizationRes{}
	}

	res := ShowOrganizationRes{
		UUID:          organization.UUID,
		Name:          organization.Name,
		Code:          organization.Code,
		Address:       organization.Address,
		Type:          organization.Type,
		Organizations: []ShowOrganizationRes{}, // Initialize as empty array
		CreatedAt:     organization.CreatedAt,
		CreatedBy:     organization.CreatedBy,
	}

	// Include parent information (without children to avoid circular reference)
	if organization.Parent != nil {
		res.Parent = &ParentOrganizationRes{
			UUID:    organization.Parent.UUID,
			Name:    organization.Parent.Name,
			Code:    organization.Parent.Code,
			Address: organization.Parent.Address,
			Type:    organization.Parent.Type,
		}
	}

	// Include children (without parent to avoid circular reference)
	if organization.Children != nil {
		for _, child := range organization.Children {
			if child != nil {
				res.Organizations = append(res.Organizations, newChildOrganizationRes(child))
			}
		}
	}

	return res
}

// newChildOrganizationRes - creates organization response without parent field to avoid circular reference
func newChildOrganizationRes(organization *entities.Organization) ShowOrganizationRes {
	if organization == nil {
		return ShowOrganizationRes{}
	}

	res := ShowOrganizationRes{
		UUID:          organization.UUID,
		Name:          organization.Name,
		Code:          organization.Code,
		Address:       organization.Address,
		Type:          organization.Type,
		Organizations: []ShowOrganizationRes{}, // Initialize as empty array
		CreatedAt:     organization.CreatedAt,
		CreatedBy:     organization.CreatedBy,
		// Parent is intentionally NOT set to avoid circular reference
	}

	// Recursively add children (but without parent)
	if organization.Children != nil {
		for _, child := range organization.Children {
			if child != nil {
				res.Organizations = append(res.Organizations, newChildOrganizationRes(child))
			}
		}
	}

	return res
}
