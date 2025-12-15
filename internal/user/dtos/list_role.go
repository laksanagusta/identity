package dtos

import (
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/nullable"
)

type PermissionResp struct {
	UUID        string              `json:"id"`
	Name        nullable.NullString `json:"name"`
	Action      nullable.NullString `json:"action"`
	Resource    nullable.NullString `json:"resource"`
	Description nullable.NullString `json:"description"`
}

type ListRoleResp struct {
	UUID        string              `json:"id"`
	Name        nullable.NullString `json:"name"`
	Description nullable.NullString `json:"description"`
	Permissions []PermissionResp    `json:"permissions"`
}

func NewListRoleResp(role []entities.Role) []ListRoleResp {
	var resp []ListRoleResp
	for _, v := range role {
		permissions := make([]PermissionResp, 0, len(v.Permissions))
		for _, p := range v.Permissions {
			permissions = append(permissions, PermissionResp{
				UUID:        p.UUID,
				Name:        p.Name,
				Action:      p.Action,
				Resource:    p.Resource,
				Description: p.Description,
			})
		}
		resp = append(resp, ListRoleResp{
			UUID:        v.UUID,
			Name:        v.Name,
			Description: v.Description,
			Permissions: permissions,
		})
	}

	return resp
}

// NewListRoleResp2 creates list role response from pointer slice (for pagination)
func NewListRoleResp2(roles []*entities.Role) []ListRoleResp {
	resp := make([]ListRoleResp, 0, len(roles))
	for _, v := range roles {
		if v != nil {
			permissions := make([]PermissionResp, 0, len(v.Permissions))
			for _, p := range v.Permissions {
				permissions = append(permissions, PermissionResp{
					UUID:        p.UUID,
					Name:        p.Name,
					Action:      p.Action,
					Resource:    p.Resource,
					Description: p.Description,
				})
			}
			resp = append(resp, ListRoleResp{
				UUID:        v.UUID,
				Name:        v.Name,
				Description: v.Description,
				Permissions: permissions,
			})
		}
	}

	return resp
}

// NewShowRoleResp creates single role response with permissions
func NewShowRoleResp(role *entities.Role) *ListRoleResp {
	if role == nil {
		return nil
	}

	permissions := make([]PermissionResp, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		permissions = append(permissions, PermissionResp{
			UUID:        p.UUID,
			Name:        p.Name,
			Action:      p.Action,
			Resource:    p.Resource,
			Description: p.Description,
		})
	}

	return &ListRoleResp{
		UUID:        role.UUID,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
	}
}
