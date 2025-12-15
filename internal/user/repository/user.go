package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/user"
	"github.com/laksanagusta/identity/pkg/database"
	"github.com/laksanagusta/identity/pkg/pagination"

	"github.com/jmoiron/sqlx"
)

func NewUserRepo(db *sqlx.DB) user.Repository {
	return &userRepo{
		conn: db,
		db:   db,
	}
}

type userRepo struct {
	conn *sqlx.DB
	db   database.Queryer
}

func (r *userRepo) Insert(ctx context.Context, user entities.User) (string, error) {
	var returnedUUID string
	err := r.db.GetContext(ctx,
		&returnedUUID,
		insertUser,
		user.UUID,
		user.EmployeeID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.PasswordHash,
		user.OrganizationUUID,
		time.Now(),
		user.CreatedBy,
		time.Now(),
		user.UpdatedBy,
	)
	if err != nil {
		return "", err
	}

	return returnedUUID, nil
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	row := r.db.QueryRowxContext(ctx, findByUsername, username)
	err := row.Scan(
		&user.UUID,
		&user.EmployeeID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.PasswordHash,
		&user.OrganizationUUID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*entities.User, error) {
	var user entities.User
	row := r.db.QueryRowxContext(ctx, findByPhoneNumber, phoneNumber)
	err := row.Scan(
		&user.UUID,
		&user.EmployeeID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepo) Update(ctx context.Context, user entities.User) error {
	res, err := r.db.ExecContext(ctx,
		updateUser,
		user.EmployeeID.IsExists,
		user.EmployeeID,
		user.PhoneNumber.IsExists,
		user.PhoneNumber,
		user.FirstName.IsExists,
		user.FirstName,
		user.LastName.IsExists,
		user.LastName,
		user.PasswordHash.IsExists,
		user.PasswordHash,
		user.CreatedBy,
		time.Now(),
		user.Username.IsExists,
		user.Username,
		user.UUID,
	)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return nil
	}

	return nil
}

func (r *userRepo) FindByUUID(ctx context.Context, uuid string) (*entities.User, error) {
	var user entities.User
	row := r.db.QueryRowxContext(ctx, findUserById, uuid)
	err := row.Scan(
		&user.UUID,
		&user.EmployeeID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.OrganizationUUID,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindByEmployeeID(ctx context.Context, employeeID string) (*entities.User, error) {
	var user entities.User
	row := r.db.QueryRowxContext(ctx, findByEmployeeId, employeeID)
	err := row.Scan(
		&user.UUID,
		&user.EmployeeID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.OrganizationUUID,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepo) FindRoleByUUID(ctx context.Context, uuid string) (*entities.Role, error) {
	var role entities.Role
	row := r.db.QueryRowxContext(ctx, findRoleById, uuid)
	err := row.Scan(
		&role.UUID,
		&role.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &role, nil
}

func (r *userRepo) FindRoleByUserUUID(ctx context.Context, uuid string) ([]*entities.Role, error) {
	rows, err := r.db.QueryxContext(ctx, findRoleByUserId, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entities.Role
	for rows.Next() {
		var role entities.Role
		err := rows.Scan(
			&role.UUID,
			&role.Name,
			&role.Description,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *userRepo) FindRole(ctx context.Context) ([]entities.Role, error) {
	var roles []entities.Role
	row, err := r.db.QueryxContext(ctx, findRole)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var role entities.Role
		err := row.Scan(
			&role.UUID,
			&role.Name,
			&role.Description,
		)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (r *userRepo) Index(ctx context.Context, params *pagination.QueryParams) ([]*entities.User, int64, error) {
	// Build count query
	countBuilder := pagination.NewQueryBuilder("SELECT COUNT(*) FROM users")
	for _, filter := range params.Filters {
		if err := countBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}
	if err := countBuilder.AddSearch(params.Search, []string{"username", "first_name"}); err != nil {
		return nil, 0, err
	}
	countQuery, countArgs := countBuilder.Build()

	var totalCount int64
	err := r.db.GetContext(ctx, &totalCount, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	// Build main query
	queryBuilder := pagination.NewQueryBuilder("SELECT * FROM users")
	for _, filter := range params.Filters {
		if err := queryBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}
	if err := queryBuilder.AddSearch(params.Search, []string{"username", "first_name"}); err != nil {
		return nil, 0, err
	}
	for _, sort := range params.Sorts {
		if err := queryBuilder.AddSort(sort); err != nil {
			return nil, 0, err
		}
	}

	query, args := queryBuilder.Build()

	// Add pagination
	offset := (params.Pagination.Page - 1) * params.Pagination.Limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Pagination.Limit, offset)

	var users []*entities.User
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var user entities.User
		if err := rows.StructScan(&user); err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

func (r *userRepo) Delete(ctx context.Context, uuid string, username string) error {
	res, err := r.db.ExecContext(ctx,
		deleteUser,
		time.Now(),
		username,
		uuid,
	)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return nil
	}

	return nil
}

func (r *userRepo) InsertUserRole(ctx context.Context, userRole entities.UserRole) (string, error) {
	var returnedUUID string
	err := r.db.GetContext(ctx,
		&returnedUUID,
		insertUserRole,
		userRole.UUID,
		userRole.UserUUID,
		userRole.RoleUUID,
		userRole.CreatedAt,
		userRole.CreatedBy,
		userRole.UpdatedAt,
		userRole.UpdatedBy,
	)
	if err != nil {
		return "", err
	}

	return returnedUUID, nil
}

func (r *userRepo) BulkInsertUserRoles(ctx context.Context, userRoles []entities.UserRole) error {
	if len(userRoles) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(userRoles))
	valueArgs := make([]interface{}, 0, len(userRoles)*7)
	for i, ur := range userRoles {
		base := i*7 + 1
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				base, base+1, base+2, base+3, base+4, base+5, base+6,
			))
		valueArgs = append(valueArgs,
			ur.UUID,
			ur.UserUUID,
			ur.RoleUUID,
			ur.CreatedAt,
			ur.CreatedBy,
			ur.UpdatedAt,
			ur.UpdatedBy,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO user_roles (uuid, user_uuid, role_uuid, created_at, created_by, updated_at, updated_by)
		VALUES %s
	`, strings.Join(valueStrings, ","))

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) DeleteUserRole(ctx context.Context, uuid string) error {
	res, err := r.db.ExecContext(ctx,
		deleteUserRole,
		uuid,
	)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return nil
	}

	return nil
}

func (r *userRepo) DeleteUserRoleByUserUUID(ctx context.Context, userUUID string) error {
	res, err := r.db.ExecContext(ctx,
		deleteUserRoleByUserUUID,
		userUUID,
	)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return nil
	}

	return nil
}

func (r *userRepo) InsertRole(ctx context.Context, role entities.Role) (string, error) {
	var returnedUUID string
	err := r.db.GetContext(ctx,
		&returnedUUID,
		insertRole,
		role.UUID,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.CreatedBy,
		role.UpdatedAt,
		role.UpdatedBy,
	)
	if err != nil {
		return "", err
	}

	return returnedUUID, nil
}

func (r *userRepo) DeleteRole(ctx context.Context, uuid string) error {
	res, err := r.db.ExecContext(ctx,
		deleteRole,
		uuid,
	)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return nil
	}

	return nil
}

func (r *userRepo) FindRoleWithPermissions(ctx context.Context, uuid string) (*entities.Role, error) {
	role, err := r.FindRoleByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	permissions, err := r.FindPermissionByRoleUUIDs(ctx, []string{uuid})
	if err != nil {
		return nil, err
	}

	// Convert []*Permission to []Permission for entity
	for _, p := range permissions {
		if p != nil {
			role.Permissions = append(role.Permissions, *p)
		}
	}

	return role, nil
}

func (r *userRepo) UpdateRole(ctx context.Context, role entities.Role) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE roles SET
			name = CASE WHEN $1 THEN $2 ELSE name END,
			description = CASE WHEN $3 THEN $4 ELSE description END,
			updated_at = $5,
			updated_by = $6
		WHERE uuid = $7`,
		role.Name.IsExists,
		role.Name,
		role.Description.IsExists,
		role.Description,
		time.Now(),
		role.UpdatedBy,
		role.UUID,
	)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return errors.New("no row affected")
	}

	return nil
}

func (r *userRepo) BulkInsertRolePermissions(ctx context.Context, rolePermissions []entities.RolaPermission) error {
	if len(rolePermissions) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(rolePermissions))
	valueArgs := make([]interface{}, 0, len(rolePermissions)*7)
	for i, rp := range rolePermissions {
		base := i*7 + 1
		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				base, base+1, base+2, base+3, base+4, base+5, base+6,
			))
		valueArgs = append(valueArgs,
			rp.UUID,
			rp.RoleUUID,
			rp.PermissionUUID,
			rp.CreatedAt,
			rp.CreatedBy,
			rp.UpdatedAt,
			rp.UpdatedBy,
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO role_permissions (uuid, role_uuid, permission_uuid, created_at, created_by, updated_at, updated_by)
		VALUES %s
	`, strings.Join(valueStrings, ","))

	_, err := r.db.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) DeleteRolePermissionsByRoleUUID(ctx context.Context, roleUUID string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM role_permissions WHERE role_uuid = $1`,
		roleUUID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) FindUserRoleByUUID(ctx context.Context, uuid string) (*entities.UserRole, error) {
	var userRole entities.UserRole
	row := r.db.QueryRowxContext(ctx, findUserRoleById, uuid)
	err := row.Scan(
		&userRole.UUID,
		&userRole.UserUUID,
		&userRole.RoleUUID,
		&userRole.CreatedAt,
		&userRole.CreatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &userRole, nil
}

func (r *userRepo) FindRoleByName(ctx context.Context, name string) (*entities.Role, error) {
	var role entities.Role
	row := r.db.QueryRowxContext(ctx, findRoleByName, name)
	err := row.Scan(
		&role.UUID,
		&role.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &role, nil
}

func (r *userRepo) InsertPermission(ctx context.Context, permission entities.Permission) (string, error) {
	var returnedUUID string
	err := r.db.GetContext(ctx,
		&returnedUUID,
		insertPermission,
		permission.UUID,
		permission.Name,
		permission.Action,
		permission.Resource,
		permission.CreatedAt,
		permission.CreatedBy,
		permission.UpdatedAt,
		permission.UpdatedBy,
	)
	if err != nil {
		return "", err
	}

	return returnedUUID, nil
}

func (r *userRepo) DeletePermission(ctx context.Context, uuid string) error {
	res, err := r.db.ExecContext(ctx,
		deletePermission,
		uuid,
	)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return nil
	}

	return nil
}

func (r *userRepo) FindPermissionByUUID(ctx context.Context, uuid string) (*entities.Permission, error) {
	var permission entities.Permission
	row := r.db.QueryRowxContext(ctx, findPermissionById, uuid)
	err := row.Scan(
		&permission.UUID,
		&permission.Name,
		&permission.Action,
		&permission.Resource,
		&permission.Description,
		&permission.CreatedAt,
		&permission.CreatedBy,
		&permission.UpdatedAt,
		&permission.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &permission, nil
}

func (r *userRepo) UpdatePermission(ctx context.Context, permission entities.Permission) error {
	res, err := r.db.ExecContext(ctx,
		updatePermission,
		permission.Name.IsExists,
		permission.Name,
		permission.Action.IsExists,
		permission.Action,
		permission.Resource.IsExists,
		permission.Resource,
		permission.Description.IsExists,
		permission.Description,
		permission.UpdatedBy,
		permission.UpdatedAt,
		permission.UUID,
	)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return nil
	}

	return nil
}

func (r *userRepo) FindSamePermission(ctx context.Context, permission entities.Permission) (*entities.Permission, error) {
	row := r.db.QueryRowxContext(ctx, findSamePermission, permission.Name, permission.Action, permission.Resource)
	err := row.Scan(
		&permission.UUID,
		&permission.Name,
		&permission.Action,
		&permission.Resource,
		&permission.Description,
		&permission.CreatedAt,
		&permission.CreatedBy,
		&permission.UpdatedAt,
		&permission.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &permission, nil
}

func (r *userRepo) FindSamePermissionExcludeCurrent(ctx context.Context, permission entities.Permission) (*entities.Permission, error) {
	row := r.db.QueryRowxContext(ctx, findSamePermissionExcludeCurrent, permission.Action, permission.Resource, permission.UUID)
	err := row.Scan(
		&permission.UUID,
		&permission.Name,
		&permission.Action,
		&permission.Resource,
		&permission.Description,
		&permission.CreatedAt,
		&permission.CreatedBy,
		&permission.UpdatedAt,
		&permission.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &permission, nil
}

func (r *userRepo) InsertRolePermission(ctx context.Context, rolePermission entities.RolaPermission) (string, error) {
	var returnedUUID string
	err := r.db.GetContext(ctx,
		&returnedUUID,
		insertRolePermission,
		rolePermission.UUID,
		rolePermission.RoleUUID,
		rolePermission.PermissionUUID,
		rolePermission.CreatedAt,
		rolePermission.CreatedBy,
		rolePermission.UpdatedAt,
		rolePermission.UpdatedBy,
	)
	if err != nil {
		return "", err
	}

	return returnedUUID, nil
}

func (r *userRepo) DeleteRolePermission(ctx context.Context, uuid string) error {
	res, err := r.db.ExecContext(ctx,
		deleteRolePermission,
		uuid,
	)
	if err != nil {
		return err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowAffected == 0 {
		return nil
	}

	return nil
}

func (r *userRepo) FindRolePermissionByUUID(ctx context.Context, uuid string) (*entities.RolaPermission, error) {
	var rolePermission entities.RolaPermission
	row := r.db.QueryRowxContext(ctx, findRolePermissionById, uuid)
	err := row.Scan(
		&rolePermission.UUID,
		&rolePermission.RoleUUID,
		&rolePermission.PermissionUUID,
		&rolePermission.CreatedAt,
		&rolePermission.CreatedBy,
		&rolePermission.UpdatedAt,
		&rolePermission.UpdatedBy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rolePermission, nil
}

func (r *userRepo) FindPermissionByRoleUUIDs(ctx context.Context, roleUUIDs []string) ([]*entities.Permission, error) {
	if len(roleUUIDs) == 0 {
		return []*entities.Permission{}, nil
	}

	query, args, err := sqlx.In(findPermissionByRoleUUIDs, roleUUIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*entities.Permission
	for rows.Next() {
		var perm entities.Permission
		err := rows.Scan(
			&perm.UUID,
			&perm.Name,
			&perm.Action,
			&perm.Resource,
			&perm.Description,
			&perm.CreatedAt,
			&perm.CreatedBy,
			&perm.UpdatedAt,
			&perm.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &perm)
	}

	return permissions, nil
}

func (r *userRepo) FindUserRolesByUserUUIDs(ctx context.Context, userUUIDs []string) ([]*entities.UserRole, error) {
	if len(userUUIDs) == 0 {
		return []*entities.UserRole{}, nil
	}

	query, args, err := sqlx.In(findUserRoleByUserUUIDs, userUUIDs)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userRoles []*entities.UserRole
	for rows.Next() {
		var (
			userRole entities.UserRole
			role     entities.Role
		)

		err := rows.Scan(
			&userRole.UUID,
			&userRole.UserUUID,
			&userRole.RoleUUID,
			&role.UUID,
			&role.Name,
			&role.Description,
			&userRole.CreatedAt,
			&userRole.CreatedBy,
			&userRole.UpdatedAt,
			&userRole.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		userRole.Role = &role
		userRoles = append(userRoles, &userRole)
	}

	return userRoles, nil
}

func (r *userRepo) IndexPermission(ctx context.Context, params *pagination.QueryParams) ([]*entities.Permission, int64, error) {
	countBuilder := pagination.NewQueryBuilder("SELECT COUNT(*) FROM permissions")
	for _, filter := range params.Filters {
		if err := countBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}
	if err := countBuilder.AddSearch(params.Search, []string{"name"}); err != nil {
		return nil, 0, err
	}
	countQuery, countArgs := countBuilder.Build()

	var totalCount int64
	err := r.db.GetContext(ctx, &totalCount, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	queryBuilder := pagination.NewQueryBuilder("SELECT * FROM permissions")
	for _, filter := range params.Filters {
		if err := queryBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}
	if err := queryBuilder.AddSearch(params.Search, []string{"name"}); err != nil {
		return nil, 0, err
	}
	for _, sort := range params.Sorts {
		if err := queryBuilder.AddSort(sort); err != nil {
			return nil, 0, err
		}
	}

	query, args := queryBuilder.Build()

	offset := (params.Pagination.Page - 1) * params.Pagination.Limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Pagination.Limit, offset)

	var permissions []*entities.Permission
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var permission entities.Permission
		if err := rows.StructScan(&permission); err != nil {
			return nil, 0, err
		}
		permissions = append(permissions, &permission)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return permissions, totalCount, nil
}

func (r *userRepo) IndexRole(ctx context.Context, params *pagination.QueryParams) ([]*entities.Role, int64, error) {
	countBuilder := pagination.NewQueryBuilder("SELECT COUNT(*) FROM roles")
	for _, filter := range params.Filters {
		if err := countBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}
	if err := countBuilder.AddSearch(params.Search, []string{"name"}); err != nil {
		return nil, 0, err
	}
	countQuery, countArgs := countBuilder.Build()

	var totalCount int64
	err := r.db.GetContext(ctx, &totalCount, countQuery, countArgs...)
	if err != nil {
		return nil, 0, err
	}

	queryBuilder := pagination.NewQueryBuilder("SELECT * FROM roles")
	for _, filter := range params.Filters {
		if err := queryBuilder.AddFilter(filter); err != nil {
			return nil, 0, err
		}
	}
	if err := queryBuilder.AddSearch(params.Search, []string{"name"}); err != nil {
		return nil, 0, err
	}
	for _, sort := range params.Sorts {
		if err := queryBuilder.AddSort(sort); err != nil {
			return nil, 0, err
		}
	}

	query, args := queryBuilder.Build()

	offset := (params.Pagination.Page - 1) * params.Pagination.Limit
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", params.Pagination.Limit, offset)

	var roles []*entities.Role
	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var role entities.Role
		if err := rows.StructScan(&role); err != nil {
			return nil, 0, err
		}
		roles = append(roles, &role)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return roles, totalCount, nil
}
