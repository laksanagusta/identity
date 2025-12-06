package external

import (
	"time"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"
)

// ExternalOrganizationRes response structure untuk external API
type ExternalOrganizationRes struct {
	UUID      string              `json:"id"`
	Name      nullable.NullString `json:"name"`
	Code      nullable.NullString `json:"code"`
	Address   nullable.NullString `json:"address"`
	Type      nullable.NullString `json:"type"`
	ParentUUID nullable.NullString `json:"parent_id"`
	Level     nullable.NullInt32  `json:"level"`
	Path      nullable.NullString `json:"path"`
	IsActive  bool                `json:"is_active"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

// NewExternalOrganizationRes convert dari entities.Organization ke ExternalOrganizationRes
func NewExternalOrganizationRes(organization entities.Organization) ExternalOrganizationRes {
	return ExternalOrganizationRes{
		UUID:       organization.UUID,
		Name:       organization.Name,
		Code:       organization.Code,
		Address:    organization.Address,
		Type:       organization.Type,
		ParentUUID: organization.ParentUUID,
		Level:      organization.Level,
		Path:       organization.Path,
		IsActive:   organization.IsActive,
		CreatedAt:  organization.CreatedAt,
		UpdatedAt:  organization.UpdatedAt,
	}
}