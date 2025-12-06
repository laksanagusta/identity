package entities

import (
	"time"

	"github.com/google/uuid"
)

type BaseModel struct {
	UUID      string    `json:"id" db:"uuid"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy string    `json:"updated_by" db:"updated_by"`
}

type SoftDeleteModel struct {
	BaseModel
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	DeletedBy *string    `json:"deleted_by,omitempty" db:"deleted_by"`
}

func NewBaseModel(username string) BaseModel {
	now := time.Now()
	return BaseModel{
		UUID:      uuid.New().String(),
		CreatedAt: now,
		CreatedBy: username,
		UpdatedAt: now,
		UpdatedBy: username,
	}
}

func (m *BaseModel) UpdateModel(username string) {
	m.UpdatedAt = time.Now()
	m.UpdatedBy = username
}
