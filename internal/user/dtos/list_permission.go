package dtos

import (
	"time"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"
)

type ListPermissionRespData struct {
	UUID      string              `json:"id"`
	Name      nullable.NullString `json:"name"`
	Action    nullable.NullString `json:"action"`
	Resource  nullable.NullString `json:"resource"`
	CreatedAt nullable.NullString `json:"created_at"`
}

type ListPermissionRespMetadata struct {
	Count       int64 `json:"count"`
	TotalCount  int64 `json:"total_count"`
	CurrentPage int64 `json:"current_page"`
	TotalPage   int64 `json:"total_page"`
}

func NewListPermissionResp(permissions []*entities.Permission) []ListPermissionRespData {
	data := make([]ListPermissionRespData, len(permissions))
	for i, permission := range permissions {
		data[i] = ListPermissionRespData{
			UUID:     permission.UUID,
			Name:     permission.Name,
			Action:   permission.Action,
			Resource: permission.Resource,

			CreatedAt: nullable.NewString(permission.CreatedAt.Format(time.RFC3339)),
		}
	}

	return data
}
