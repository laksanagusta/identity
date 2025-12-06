package user

import (
	"context"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/user/dtos"
	"github.com/laksanagusta/identity/pkg/pagination"
)

type UseCase interface {
	Create(ctx context.Context, req dtos.CreateNewUserReq) (string, error)
	Update(ctx context.Context, cred entities.AuthenticatedUser, req dtos.UpdateUserReq) error
	Show(ctx context.Context, uuid string) (*entities.User, []string, error)
	Login(ctx context.Context, req dtos.LoginReq) (string, error)
	Index(ctx context.Context, params *pagination.QueryParams) ([]*entities.User, *pagination.PagedResponse, error)
	Delete(ctx context.Context, cred entities.AuthenticatedUser, uuid string) error
	ChangePassword(ctx context.Context, cred entities.AuthenticatedUser, req dtos.ChangePassword) error

	Role(ctx context.Context) ([]entities.Role, error)
	CreateRole(ctx context.Context, role entities.Role) error
	DeleteRole(ctx context.Context, uuid string) error

	CreateUserRole(ctx context.Context, userRole entities.UserRole) error
	DeleteUserRole(ctx context.Context, uuid string) error

	CreatePermission(ctx context.Context, permission entities.Permission) error
	DeletePermission(ctx context.Context, uuid string) error
	UpdatePermission(ctx context.Context, permission entities.Permission) error
	IndexPermission(ctx context.Context, params *pagination.QueryParams) ([]*entities.Permission, *pagination.PagedResponse, error)

	CreateRolePermission(ctx context.Context, rolePermission entities.RolaPermission) error
	DeleteRolePermission(ctx context.Context, uuid string) error
}
