package usecase

import (
	"context"

	"github.com/laksanagusta/identity/constants"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/organization"
	"github.com/laksanagusta/identity/internal/organization/dtos"
	"github.com/laksanagusta/identity/internal/user"
	"github.com/laksanagusta/identity/pkg/database"
	"github.com/laksanagusta/identity/pkg/errorhelper"
)

type UseCaseParameter struct {
	OrganizationRepo organization.Repository
	TxManager        database.Manager
	UserUC           user.UseCase
}

func NewOrganizationUseCase(uc UseCaseParameter) organization.UseCase {
	return &OrganizationUseCase{
		organizationRepo: uc.OrganizationRepo,
		txManager:        uc.TxManager,
		userUC:           uc.UserUC,
	}
}

type OrganizationUseCase struct {
	organizationRepo organization.Repository
	txManager        database.Manager
	userUC           user.UseCase
}

func (uc *OrganizationUseCase) Create(ctx context.Context, cred entities.AuthenticatedUser, req dtos.CreateNewOrganizationReq) (string, error) {
	organization := req.NewOrganization(cred)

	var newOrganizationUUID string
	err := uc.txManager.Atomic(ctx, func(ctx context.Context, tx database.DBTx) error {
		organizationRepoTrx := uc.organizationRepo.WithTransaction(tx)

		var parentPath string
		if organization.ParentUUID.IsExists {
			parent, err := organizationRepoTrx.FindOrganizationByUUID(ctx, *organization.ParentUUID.Val)
			if err != nil {
				return err
			}

			parentPath = parent.Path.GetOrDefault()
		}

		organization.BuildPath(parentPath)

		newUUID, err := organizationRepoTrx.Insert(ctx, organization)
		if err != nil {
			return err
		}
		newOrganizationUUID = newUUID
		return nil
	})
	if err != nil {
		return "", err
	}

	return newOrganizationUUID, nil
}

func (uc *OrganizationUseCase) Update(ctx context.Context, cred entities.AuthenticatedUser, req dtos.UpdateOrganizationReq) error {
	organization := req.NewUpdateOrganization(cred)

	existingOrganization, err := uc.organizationRepo.FindOrganizationByUUID(ctx, req.OrganizationUUID)
	if err != nil {
		return err
	}
	if existingOrganization == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"organization_id": {constants.ErrMsgNotFound},
		})
	}

	err = uc.organizationRepo.Update(ctx, organization)
	if err != nil {
		return err
	}

	return nil
}

func (uc *OrganizationUseCase) ListOrganization(ctx context.Context, cred entities.AuthenticatedUser, req dtos.ListOrganizationReq) ([]entities.Organization, *entities.Metadata, error) {
	listOrganizationParams, err := req.NewListOrganizationParams()
	if err != nil {
		return nil, nil, err
	}

	organizations, metadata, err := uc.organizationRepo.IndexOrganization(ctx, listOrganizationParams)
	if err != nil {
		return nil, nil, err
	}

	return organizations, metadata, nil
}

func (uc *OrganizationUseCase) Show(ctx context.Context, cred entities.AuthenticatedUser, uuid string) (*entities.Organization, error) {
	existingOrganization, err := uc.organizationRepo.FindOrganizationByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if existingOrganization == nil {
		return nil, errorhelper.BadRequestMap(map[string][]string{
			"organization_id": {constants.ErrMsgNotFound},
		})
	}

	return existingOrganization, nil
}

func (uc *OrganizationUseCase) Delete(ctx context.Context, cred entities.AuthenticatedUser, uuid string) error {
	organization, err := uc.organizationRepo.FindOrganizationByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	if organization == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"organization_id": {constants.ErrMsgNotFound},
		})
	}

	err = uc.organizationRepo.Delete(ctx, uuid, cred.Username)
	if err != nil {
		return err
	}

	return nil
}
