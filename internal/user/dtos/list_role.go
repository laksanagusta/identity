package dtos

import (
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"
)

type ListRoleResp struct {
	UUID        string              `json:"id"`
	Name        nullable.NullString `json:"name"`
	Description nullable.NullString `json:"description"`
}

func NewListRoleResp(role []entities.Role) []ListRoleResp {
	var resp []ListRoleResp
	for _, v := range role {
		resp = append(resp, ListRoleResp{
			UUID:        v.UUID,
			Name:        v.Name,
			Description: v.Description,
		})
	}

	return resp
}
