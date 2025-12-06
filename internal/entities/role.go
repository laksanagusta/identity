package entities

import "github.com/laksanagusta/identity/pkg/nullable"

type Role struct {
	SoftDeleteModel
	Name        nullable.NullString `json:"name" db:"name"`
	Description nullable.NullString `json:"description" db:"description"`
	IsSystem    bool                `json:"is_system" db:"is_system"`

	Permissions []Permission `json:"permissions,omitempty" db:"-"`
}

type RolaPermission struct {
	BaseModel
	RoleUUID       string `json:"role_id" db:"role_uuid"`
	PermissionUUID string `json:"permission_id" db:"permission_uuid"`

	Role       *Role       `json:"role,omitempty" db:"-"`
	Permission *Permission `json:"permission,omitempty" db:"-"`
}

type Roles []*Role

func (r Roles) Uuids() []string {
	uuids := make([]string, 0, len(r))
	for _, role := range r {
		uuids = append(uuids, role.UUID)
	}
	return uuids
}
