package dtos

import (
	"strings"
	"time"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/errorhelper"
	"github.com/laksanagusta/identity/pkg/helper"
	"github.com/laksanagusta/identity/pkg/nullable"

	"github.com/invopop/validation"
)

type ListOrganizationReq struct {
	Page      int        `query:"page"`
	Limit     int        `query:"limit"`
	Search    *string    `query:"search"`
	Sort      *string    `query:"sort"`
	StartTime *time.Time `query:"start_time"`
	EndTime   *time.Time `query:"end_time"`
}

type ListOrganizationRespData struct {
	UUID      string              `json:"id"`
	Name      string              `json:"name"`
	Address   nullable.NullString `json:"address"`
	Type      nullable.NullString `json:"type"`
	CreatedAt nullable.NullString `json:"created_at"`
	CreatedBy string              `json:"created_by"`
}

type ListOrganizationRespMetadata struct {
	Count       int64 `json:"count"`
	TotalCount  int64 `json:"total_count"`
	CurrentPage int64 `json:"current_page"`
	TotalPage   int64 `json:"total_page"`
}

type ListOrganizationResp struct {
	Data     []ListOrganizationRespData   `json:"data"`
	Metadata ListOrganizationRespMetadata `json:"metadata"`
}

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

func (r *ListOrganizationReq) NewListOrganizationParams() (entities.ListOrganizationParams, error) {
	sortableField := map[string]string{
		"name":       "s.name",
		"address":    "s.address",
		"created_at": "s.created_at",
	}

	listOrganizationParams := entities.ListOrganizationParams{
		Offset: (r.Page - 1) * r.Limit,
		Limit:  r.Limit,
	}

	if r.StartTime != nil {
		listOrganizationParams.StartTime = nullable.NewTime(*r.StartTime)
	}

	if r.EndTime != nil {
		listOrganizationParams.EndTime = nullable.NewTime(*r.EndTime)
	}

	if r.Search != nil {
		listOrganizationParams.Search = nullable.NewString(*r.Search)
	}

	if r.Sort != nil {
		fieldName, err := helper.ValidateSortV2(sortableField, *r.Sort)
		if err != nil {
			return entities.ListOrganizationParams{}, err
		}

		splitStr := strings.Split(*r.Sort, " ")
		if len(splitStr) == 2 {
			listOrganizationParams.Sort = &entities.Sort{
				FieldName: fieldName,
				SortType:  splitStr[1],
			}
		} else {
			return entities.ListOrganizationParams{}, errorhelper.BadRequestMap(map[string][]string{
				"sort": {"invalid"},
			})
		}
	}

	return listOrganizationParams, nil
}

func NewListOrganizationResp(organizations []entities.Organization, metadata *entities.Metadata) ListOrganizationResp {
	data := make([]ListOrganizationRespData, len(organizations))

	for k, organization := range organizations {
		data[k].UUID = organization.UUID
		data[k].Name = organization.Name.GetOrDefault()
		data[k].Address = organization.Address
		data[k].Type = organization.Type
		data[k].CreatedAt = nullable.NewString(organization.CreatedAt.Format("2006-01-02T15:04:05+0700"))
		data[k].CreatedBy = organization.CreatedBy
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
