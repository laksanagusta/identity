package organization

import (
	"context"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/organization/dtos"
)

type UseCase interface {
	Create(ctx context.Context, auth entities.AuthenticatedUser, req dtos.CreateNewOrganizationReq) (string, error)
	Update(ctx context.Context, cred entities.AuthenticatedUser, req dtos.UpdateOrganizationReq) error
	Show(ctx context.Context, cred entities.AuthenticatedUser, uuid string) (*entities.Organization, error)
	ListOrganization(ctx context.Context, cred entities.AuthenticatedUser, req dtos.ListOrganizationReq) ([]entities.Organization, *entities.Metadata, error)
	Delete(ctx context.Context, cred entities.AuthenticatedUser, uuid string) error
}
