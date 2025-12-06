package usecase

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/laksanagusta/identity/constants"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/organization"
	"github.com/laksanagusta/identity/internal/user"
	"github.com/laksanagusta/identity/internal/user/dtos"
	"github.com/laksanagusta/identity/pkg/authservice/jwt"
	"github.com/laksanagusta/identity/pkg/errorhelper"
	"github.com/laksanagusta/identity/pkg/helper"
	"github.com/laksanagusta/identity/pkg/nullable"
	"github.com/laksanagusta/identity/pkg/pagination"

	"golang.org/x/crypto/bcrypt"
)

type UseCaseParameter struct {
	UserRepo         user.Repository
	OrganizationRepo organization.Repository
	JwtAuth          jwt.JwtAuth
}

func NewUserUseCase(uc UseCaseParameter) user.UseCase {
	return &UserUseCase{
		userRepo:         uc.UserRepo,
		jwtAuth:          uc.JwtAuth,
		organizationRepo: uc.OrganizationRepo,
	}
}

type UserUseCase struct {
	userRepo         user.Repository
	jwtAuth          jwt.JwtAuth
	organizationRepo organization.Repository
}

func (uc *UserUseCase) Create(ctx context.Context, req dtos.CreateNewUserReq) (string, error) {
	user := req.NewUser()

	existedUser, err := uc.userRepo.FindByUsername(ctx, strings.ToLower(req.Username))
	if err != nil {
		return "", err
	}
	if existedUser != nil {
		return "", errorhelper.BadRequestMap(map[string][]string{
			"username": {constants.ErrMsgAlreadyExist},
		})
	}

	existedUser, err = uc.userRepo.FindByPhoneNumber(ctx, req.PhoneNumber)
	if err != nil {
		return "", err
	}
	if existedUser != nil {
		return "", errorhelper.BadRequestMap(map[string][]string{
			"phone_number": {constants.ErrMsgAlreadyExist},
		})
	}

	existedUser, err = uc.userRepo.FindByEmployeeID(ctx, req.EmployeeID)
	if err != nil {
		return "", err
	}
	if existedUser != nil {
		return "", errorhelper.BadRequestMap(map[string][]string{
			"employee_id": {constants.ErrMsgAlreadyExist},
		})
	}

	organization, err := uc.organizationRepo.FindOrganizationByUUID(ctx, req.OrganizationUUID)
	if err != nil {
		return "", err
	}
	if organization == nil {
		return "", errorhelper.BadRequestMap(map[string][]string{
			"organization_id": {constants.ErrMsgNotFound},
		})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	user.PasswordHash = nullable.NewString(string(passwordHash))

	newUUID, err := uc.userRepo.Insert(ctx, user)
	if err != nil {
		return "", err
	}

	roleList := entities.Roles(user.Roles)
	roleUUIDs := roleList.Uuids()

	now := time.Now()
	userRoles := make([]entities.UserRole, 0, len(roleUUIDs))
	for _, roleUUID := range roleUUIDs {
		userRoles = append(userRoles, entities.UserRole{
			BaseModel: entities.BaseModel{
				UUID:      uuid.NewString(),
				CreatedAt: now,
				CreatedBy: "admin",
				UpdatedAt: now,
				UpdatedBy: "admin",
			},
			UserUUID: newUUID,
			RoleUUID: roleUUID,
		})
	}

	err = uc.userRepo.BulkInsertUserRoles(ctx, userRoles)
	if err != nil {
		return "", err
	}

	return newUUID, nil
}

func (uc *UserUseCase) Update(ctx context.Context, cred entities.AuthenticatedUser, req dtos.UpdateUserReq) error {
	user := req.NewUser(cred)

	if req.Username.IsExists {
		foundUser, err := uc.userRepo.FindByUsername(ctx, user.Username.GetOrDefault())
		if err != nil {
			return err
		}
		if foundUser != nil {
			return errorhelper.BadRequestMap(map[string][]string{
				"username": {constants.ErrMsgAlreadyExist},
			})
		}
	}

	if req.EmployeeID.IsExists {
		foundUser, err := uc.userRepo.FindByEmployeeID(ctx, user.EmployeeID.GetOrDefault())
		if err != nil {
			return err
		}
		if foundUser != nil && foundUser.UUID != req.UserUUID {
			return errorhelper.BadRequestMap(map[string][]string{
				"employee_id": {constants.ErrMsgAlreadyExist},
			})
		}
	}

	if req.Password.IsExists {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password.GetOrDefault()), bcrypt.MinCost)
		if err != nil {
			return err
		}

		user.PasswordHash = nullable.NewString(string(passwordHash))
	}

	if len(req.RoleUUIDs) > 0 {
		err := uc.userRepo.DeleteUserRoleByUserUUID(ctx, user.UUID)
		if err != nil {
			return err
		}

		userRoles := make([]entities.UserRole, 0, len(req.RoleUUIDs))
		now := time.Now()
		for _, roleUUID := range req.RoleUUIDs {
			userRole := entities.UserRole{
				BaseModel: entities.BaseModel{
					UUID:      uuid.NewString(),
					CreatedAt: now,
					CreatedBy: cred.Username,
					UpdatedAt: now,
					UpdatedBy: cred.Username,
				},
				UserUUID: user.UUID,
				RoleUUID: roleUUID,
			}
			userRoles = append(userRoles, userRole)
		}

		err = uc.userRepo.BulkInsertUserRoles(ctx, userRoles)
		if err != nil {
			return err
		}
	}

	err := uc.userRepo.Update(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) Show(ctx context.Context, uuid string) (*entities.User, []string, error) {
	user, err := uc.userRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errorhelper.BadRequestMap(map[string][]string{
			"user_id": {constants.ErrMsgNotFound},
		})
	}

	roles, err := uc.userRepo.FindRoleByUserUUID(ctx, uuid)
	if err != nil {
		return nil, nil, err
	}

	user.Roles = roles

	roleList := entities.Roles(roles)
	roleUUIDs := roleList.Uuids()

	permissions, err := uc.userRepo.FindPermissionByRoleUUIDs(ctx, roleUUIDs)
	if err != nil {
		return nil, nil, err
	}

	user.Permissions = permissions

	permissionsStr := make([]string, len(permissions))

	for i, p := range permissions {
		if p.Resource.IsExists && p.Action.IsExists {
			permissionsStr[i] = *p.Resource.Val + ":" + *p.Action.Val
		}
	}

	organization, err := uc.organizationRepo.FindOrganizationByUUID(ctx, user.OrganizationUUID.GetOrDefault())
	if err != nil {
		return nil, nil, err
	}
	if organization == nil {
		return nil, nil, errorhelper.BadRequestMap(map[string][]string{
			"organization_id": {constants.ErrMsgNotFound},
		})
	}

	user.Organization = organization

	return user, nil, nil
}

func (uc *UserUseCase) Login(ctx context.Context, req dtos.LoginReq) (string, error) {
	user, err := uc.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errorhelper.BadRequestMap(map[string][]string{
			"username": {constants.ErrMsgNotFound},
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.GetOrDefault()), []byte(req.Password))
	if err != nil {
		return "", errorhelper.BadRequestMap(map[string][]string{
			"password": {"invalid"},
		})
	}

	roles, err := uc.userRepo.FindRoleByUserUUID(ctx, user.UUID)
	if err != nil {
		return "", err
	}
	if len(roles) == 0 {
		return "", errorhelper.BadRequestMap(map[string][]string{
			"roles": {constants.ErrMsgNotFound},
		})
	}

	user.Roles = roles

	organization, err := uc.organizationRepo.FindOrganizationByUUID(ctx, user.OrganizationUUID.GetOrDefault())
	if err != nil {
		return "", err
	}

	token, err := uc.jwtAuth.GenerateToken(*user, *organization)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uc *UserUseCase) Role(ctx context.Context) ([]entities.Role, error) {
	return uc.userRepo.FindRole(ctx)
}

func (uc *UserUseCase) Index(ctx context.Context, params *pagination.QueryParams) ([]*entities.User, *pagination.PagedResponse, error) {
	users, totalCount, err := uc.userRepo.Index(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	userList := entities.Users(users)
	userUUIDs := userList.Uuids()
	indexUser := map[string]int32{}

	organizationUnique := map[string]struct{}{}
	organizationUUIDs := make([]string, 0, len(users))
	for i, user := range users {
		indexUser[user.UUID] = int32(i)
		if user.OrganizationUUID.Val != nil {
			_, ok := organizationUnique[*user.OrganizationUUID.Val]
			if !ok {
				organizationUUIDs = append(organizationUUIDs, *user.OrganizationUUID.Val)
				organizationUnique[*user.OrganizationUUID.Val] = struct{}{}
			}
		}
	}

	userRoles, err := uc.userRepo.FindUserRolesByUserUUIDs(ctx, userUUIDs)
	if err != nil {
		return nil, nil, err
	}

	for _, userRole := range userRoles {
		val, ok := indexUser[userRole.UserUUID]
		if ok {
			users[val].Roles = append(users[val].Roles, userRole.Role)
		}
	}

	organizations, err := uc.organizationRepo.FindOrganizationByUUIDs(ctx, organizationUUIDs)
	if err != nil {
		return nil, nil, err
	}

	organizationsByUUID := helper.IndexBy(organizations, func(o *entities.Organization) string { return o.UUID })

	// mapping organization ke user
	for _, user := range users {
		if user.OrganizationUUID.Val != nil {
			org, ok := organizationsByUUID[*user.OrganizationUUID.Val]
			if ok {
				user.Organization = org
			}
		}
	}

	totalPages := int(totalCount) / params.Pagination.Limit
	if int(totalCount)%params.Pagination.Limit > 0 {
		totalPages++
	}

	return users, &pagination.PagedResponse{
		Page:       params.Pagination.Page,
		Limit:      params.Pagination.Limit,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}, nil
}

func (uc *UserUseCase) ChangePassword(ctx context.Context, cred entities.AuthenticatedUser, req dtos.ChangePassword) error {
	user, err := uc.userRepo.FindByUUID(ctx, cred.ID)
	if err != nil {
		return err
	}
	if user == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"user_id": {constants.ErrMsgNotFound},
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.GetOrDefault()), []byte(req.OldPassword))
	if err != nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"old_password": {"invalid"},
		})
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.MinCost)
	if err != nil {
		return err
	}

	updatePassword := entities.User{
		PasswordHash: nullable.NewString(string(passwordHash)),
	}

	updatePassword.UUID = req.UserUUID
	updatePassword.UpdatedAt = time.Now()
	updatePassword.UpdatedBy = cred.Username

	err = uc.userRepo.Update(ctx, updatePassword)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) Delete(ctx context.Context, cred entities.AuthenticatedUser, uuid string) error {
	user, err := uc.userRepo.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	if user == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"user_id": {constants.ErrMsgNotFound},
		})
	}

	err = uc.userRepo.Delete(ctx, uuid, cred.Username)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) CreateUserRole(ctx context.Context, userRole entities.UserRole) error {
	user, err := uc.userRepo.FindByUUID(ctx, userRole.UserUUID)
	if err != nil {
		return err
	}
	if user == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"user_id": {constants.ErrMsgNotFound},
		})
	}

	role, err := uc.userRepo.FindRoleByUUID(ctx, userRole.RoleUUID)
	if err != nil {
		return err
	}
	if role == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"role_id": {constants.ErrMsgNotFound},
		})
	}

	_, err = uc.userRepo.InsertUserRole(ctx, userRole)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) DeleteUserRole(ctx context.Context, uuid string) error {
	user, err := uc.userRepo.FindUserRoleByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	if user == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"user_role_id": {constants.ErrMsgNotFound},
		})
	}

	err = uc.userRepo.DeleteUserRole(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) CreateRole(ctx context.Context, role entities.Role) error {
	roleExist, err := uc.userRepo.FindRoleByName(ctx, *role.Name.Val)
	if err != nil {
		return err
	}
	if roleExist != nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"role_name": {constants.ErrMsgAlreadyExist},
		})
	}

	_, err = uc.userRepo.InsertRole(ctx, role)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) DeleteRole(ctx context.Context, uuid string) error {
	role, err := uc.userRepo.FindRoleByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	if role == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"role_id": {constants.ErrMsgNotFound},
		})
	}

	err = uc.userRepo.DeleteRole(ctx, uuid)
	if err != nil {
		log.Println("sddawdwad")
		return err
	}

	return nil
}

func (uc *UserUseCase) CreatePermission(ctx context.Context, permission entities.Permission) error {
	permExist, err := uc.userRepo.FindSamePermission(ctx, permission)
	if err != nil {
		return err
	}
	if permExist != nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"permission": {constants.ErrMsgAlreadyExist},
		})
	}

	_, err = uc.userRepo.InsertPermission(ctx, permission)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) DeletePermission(ctx context.Context, uuid string) error {
	permission, err := uc.userRepo.FindPermissionByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	if permission == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"permission_id": {constants.ErrMsgNotFound},
		})
	}

	err = uc.userRepo.DeletePermission(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) UpdatePermission(ctx context.Context, permission entities.Permission) error {
	// Check if the permission exists
	existingPerm, err := uc.userRepo.FindPermissionByUUID(ctx, permission.UUID)
	if err != nil {
		return err
	}
	if existingPerm == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"permission_id": {constants.ErrMsgNotFound},
		})
	}

	samePerm, err := uc.userRepo.FindSamePermissionExcludeCurrent(ctx, permission)
	if err != nil {
		return err
	}
	if samePerm != nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"permission": {constants.ErrMsgAlreadyExist},
		})
	}

	// Update the permission
	err = uc.userRepo.UpdatePermission(ctx, permission)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) CreateRolePermission(ctx context.Context, rolePermission entities.RolaPermission) error {
	// Check if role exists
	role, err := uc.userRepo.FindRoleByUUID(ctx, rolePermission.RoleUUID)
	if err != nil {
		return err
	}
	if role == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"role_id": {constants.ErrMsgNotFound},
		})
	}

	// Check if permission exists
	permission, err := uc.userRepo.FindPermissionByUUID(ctx, rolePermission.PermissionUUID)
	if err != nil {
		return err
	}
	if permission == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"permission_id": {constants.ErrMsgNotFound},
		})
	}

	_, err = uc.userRepo.InsertRolePermission(ctx, rolePermission)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) DeleteRolePermission(ctx context.Context, uuid string) error {
	rolePermission, err := uc.userRepo.FindRolePermissionByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	if rolePermission == nil {
		return errorhelper.BadRequestMap(map[string][]string{
			"role_permission_id": {constants.ErrMsgNotFound},
		})
	}

	err = uc.userRepo.DeleteRolePermission(ctx, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) IndexPermission(ctx context.Context, params *pagination.QueryParams) ([]*entities.Permission, *pagination.PagedResponse, error) {
	permissions, totalCount, err := uc.userRepo.IndexPermission(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(totalCount) / params.Pagination.Limit
	if int(totalCount)%params.Pagination.Limit > 0 {
		totalPages++
	}

	return permissions, &pagination.PagedResponse{
		Page:       params.Pagination.Page,
		Limit:      params.Pagination.Limit,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}, nil
}
