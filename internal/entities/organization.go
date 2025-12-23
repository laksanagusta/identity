package entities

import (
	"log"

	"github.com/laksanagusta/identity/pkg/nullable"
)

type Organization struct {
	SoftDeleteModel
	Name       nullable.NullString `json:"name" db:"name"`
	Code       nullable.NullString `json:"code" db:"code"`
	Address    nullable.NullString `json:"address" db:"address"`
	Latitude   nullable.NullString `json:"latitude" db:"latitude"`
	Longitude  nullable.NullString `json:"longitude" db:"longitude"`
	Type       nullable.NullString `json:"type" db:"type"`
	ParentUUID nullable.NullString `json:"parent_id" db:"parent_uuid"`
	Level      nullable.NullInt32  `json:"level" db:"level"`
	Path       nullable.NullString `json:"path" db:"path"`
	IsActive   bool                `json:"is_active" db:"is_active"`

	Parent   *Organization   `json:"parent,omitempty" db:"-"`
	Children []*Organization `json:"children,omitempty" db:"-"`
	Users    []*User         `json:"users,omitempty" db:"-"`
}

type Organizations []*Organization

func (o *Organization) BuildPath(parentPath string) {
	if parentPath == "" {
		o.Path = nullable.NewString(o.UUID)
	} else {
		log.Println("parent path", parentPath)
		o.Path = nullable.NewString(parentPath + "." + o.UUID)
	}
}

func (os Organizations) Uuids() []string {
	uuids := make([]string, 0, len(os))
	for _, o := range os {
		if o != nil {
			uuids = append(uuids, o.UUID)
		}
	}
	return uuids
}

type ListOrganizationParams struct {
	Offset    int
	Limit     int
	StartTime nullable.NullTime
	EndTime   nullable.NullTime
	Search    nullable.NullString
	Sort      *Sort
}

type ListOrganizationProductStockParams struct {
	Offset           int
	Limit            int
	StartTime        nullable.NullTime
	EndTime          nullable.NullTime
	Search           nullable.NullString
	Sort             *Sort
	OrganizationUUID nullable.NullString
}

type ListOrganizationDebtParams struct {
	Offset           int
	Limit            int
	StartTime        nullable.NullTime
	EndTime          nullable.NullTime
	Search           nullable.NullString
	Sort             *Sort
	OrganizationUUID nullable.NullString
}

type Sort struct {
	FieldName string
	SortType  string
}

type Metadata struct {
	Count       float64
	TotalCount  float64
	CurrentPage float64
	TotalPage   float64
}
