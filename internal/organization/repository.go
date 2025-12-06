package organization

import (
	"context"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/database"
)

type Repository interface {
	WithTransaction(tx database.DBTx) Repository

	Insert(ctx context.Context, organization entities.Organization) (string, error)
	FindOrganizationByUUID(ctx context.Context, uuid string) (*entities.Organization, error)
	Update(ctx context.Context, organization entities.Organization) error
	IndexOrganization(ctx context.Context, params entities.ListOrganizationParams) ([]entities.Organization, *entities.Metadata, error)
	Delete(ctx context.Context, uuid string, username string) error
	FindOrganizationByUUIDs(ctx context.Context, uuids []string) ([]*entities.Organization, error)
}
