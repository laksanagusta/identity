package user

import (
	"context"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/pkg/pagination"
)

type Repository interface {
	Insert(ctx context.Context, user entities.User) (string, error)
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error)
	FindByEmployeeID(ctx context.Context, employeeID string) (*entities.User, error)
	Update(ctx context.Context, user entities.User) error
	FindByUUID(ctx context.Context, uuid string) (*entities.User, error)
	Index(ctx context.Context, params *pagination.QueryParams) ([]*entities.User, int64, error)
	Delete(ctx context.Context, uuid string, username string) error

	// role
	FindRoleByUUID(ctx context.Context, uuid string) (*entities.Role, error)
	FindRole(ctx context.Context) ([]entities.Role, error)
	FindRoleByName(ctx context.Context, name string) (*entities.Role, error)
	DeleteRole(ctx context.Context, uuid string) error
	InsertRole(ctx context.Context, role entities.Role) (string, error)

	// user-role
	InsertUserRole(ctx context.Context, userRole entities.UserRole) (string, error)
	DeleteUserRole(ctx context.Context, uuid string) error
	FindUserRoleByUUID(ctx context.Context, uuid string) (*entities.UserRole, error)
	FindRoleByUserUUID(ctx context.Context, uuid string) ([]*entities.Role, error)
	BulkInsertUserRoles(ctx context.Context, userRoles []entities.UserRole) error
	FindUserRolesByUserUUIDs(ctx context.Context, userUUIDs []string) ([]*entities.UserRole, error)
	DeleteUserRoleByUserUUID(ctx context.Context, userUUID string) error

	// permission
	InsertPermission(ctx context.Context, permission entities.Permission) (string, error)
	DeletePermission(ctx context.Context, uuid string) error
	FindPermissionByUUID(ctx context.Context, uuid string) (*entities.Permission, error)
	UpdatePermission(ctx context.Context, permission entities.Permission) error
	FindSamePermission(ctx context.Context, permission entities.Permission) (*entities.Permission, error)
	FindSamePermissionExcludeCurrent(ctx context.Context, permission entities.Permission) (*entities.Permission, error)
	IndexPermission(ctx context.Context, params *pagination.QueryParams) ([]*entities.Permission, int64, error)

	// role-permission
	InsertRolePermission(ctx context.Context, rolePermission entities.RolaPermission) (string, error)
	DeleteRolePermission(ctx context.Context, uuid string) error
	FindRolePermissionByUUID(ctx context.Context, uuid string) (*entities.RolaPermission, error)
	FindPermissionByRoleUUIDs(ctx context.Context, roleUUIDs []string) ([]*entities.Permission, error)
}
