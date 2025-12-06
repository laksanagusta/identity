package entities

import "github.com/laksanagusta/identity/pkg/nullable"

// Permission dengan action dan resource
type Permission struct {
	BaseModel
	Name        nullable.NullString `json:"name" db:"name"`
	Action      nullable.NullString `json:"action" db:"action"`     // create, read, update, delete, approve
	Resource    nullable.NullString `json:"resource" db:"resource"` // proposal, workflow, user
	Description nullable.NullString `json:"description" db:"description"`
}

type Permissions []*Permission

func (ps Permissions) Uuids() []string {
	uuids := make([]string, 0, len(ps))
	for _, p := range ps {
		if p != nil {
			uuids = append(uuids, p.UUID)
		}
	}
	return uuids
}
