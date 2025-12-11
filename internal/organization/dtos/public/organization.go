package public

import (
	"time"

	"github.com/invopop/validation"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/organization/dtos"
	"github.com/laksanagusta/identity/pkg/nullable"
)

// ListOrganizationReq represents the public API request for listing organizations
type ListOrganizationReq struct {
	Page      int        `query:"page"`
	Limit     int        `query:"limit"`
	Search    *string    `query:"search"`
	Sort      *string    `query:"sort"`
	StartTime *time.Time `query:"start_time"`
	EndTime   *time.Time `query:"end_time"`
}

// ListOrganizationRespData represents a single organization in the public API response
type ListOrganizationRespData struct {
	UUID      string              `json:"id"`
	Name      string              `json:"name"`
	Address   nullable.NullString `json:"address"`
	Type      nullable.NullString `json:"type"`
	CreatedAt nullable.NullString `json:"created_at"`
}

// ListOrganizationRespMetadata represents pagination metadata
type ListOrganizationRespMetadata struct {
	Count       int64 `json:"count"`
	TotalCount  int64 `json:"total_count"`
	CurrentPage int64 `json:"current_page"`
	TotalPage   int64 `json:"total_page"`
}

// ListOrganizationResp represents the full public API response for listing organizations
type ListOrganizationResp struct {
	Data     []ListOrganizationRespData   `json:"data"`
	Metadata ListOrganizationRespMetadata `json:"metadata"`
}

// PublicOrganizationRes represents a single organization in public API response
type PublicOrganizationRes struct {
	UUID      string              `json:"id"`
	Name      string              `json:"name"`
	Address   nullable.NullString `json:"address"`
	Type      nullable.NullString `json:"type"`
	CreatedAt nullable.NullString `json:"created_at"`
}

// Validate validates the list organization request
func (r ListOrganizationReq) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Limit, validation.Required, validation.Min(1)),
		validation.Field(&r.Page, validation.Required, validation.Min(1)),
		validation.Field(&r.Search, validation.Length(1, 255)),
		validation.Field(&r.Sort, validation.Length(1, 255)),
		validation.Field(&r.StartTime, validation.Length(1, 255), validation.Date(time.RFC3339)),
		validation.Field(&r.EndTime, validation.Length(1, 255), validation.Date(time.RFC3339)),
	)

	return err
}

// ToInternalReq converts public request to internal request
func (r *ListOrganizationReq) ToInternalReq() dtos.ListOrganizationReq {
	return dtos.ListOrganizationReq{
		Page:      r.Page,
		Limit:     r.Limit,
		Search:    r.Search,
		Sort:      r.Sort,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
	}
}

// NewListOrganizationResp creates a new public list organization response
func NewListOrganizationResp(organizations []entities.Organization, metadata *entities.Metadata) ListOrganizationResp {
	data := make([]ListOrganizationRespData, len(organizations))

	for k, organization := range organizations {
		data[k].UUID = organization.UUID
		data[k].Name = organization.Name.GetOrDefault()
		data[k].Address = organization.Address
		data[k].Type = organization.Type
		data[k].CreatedAt = nullable.NewString(organization.CreatedAt.Format("2006-01-02T15:04:05+0700"))
	}

	return ListOrganizationResp{
		Data: data,
		Metadata: ListOrganizationRespMetadata{
			Count:       int64(metadata.Count),
			TotalCount:  int64(metadata.TotalCount),
			CurrentPage: int64(metadata.CurrentPage),
			TotalPage:   int64(metadata.TotalPage),
		},
	}
}

// NewPublicOrganizationRes creates a new public organization response
func NewPublicOrganizationRes(organization entities.Organization) PublicOrganizationRes {
	return PublicOrganizationRes{
		UUID:      organization.UUID,
		Name:      organization.Name.GetOrDefault(),
		Address:   organization.Address,
		Type:      organization.Type,
		CreatedAt: nullable.NewString(organization.CreatedAt.Format("2006-01-02T15:04:05+0700")),
	}
}
