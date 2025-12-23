package repository

var (
	insertUser = `INSERT INTO users (
		uuid,
		employee_id,
		username,
		first_name,
		last_name,
		phone_number,
		password_hash,
		organization_uuid,
		is_approved,
		avatar_gradient_start,
		avatar_gradient_end,
		created_at,
		created_by,
		updated_at,
		updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING uuid`

	findByUsername = `SELECT
		uuid,
		employee_id,
		username,
		first_name,
		last_name,
		phone_number,
		password_hash,
		organization_uuid,
		is_approved
	FROM users
	WHERE username = $1 AND deleted_at is null LIMIT 1`

	findByPhoneNumber = `SELECT
		uuid,
		employee_id,
		username,
		first_name,
		last_name,
		phone_number
	FROM users
	WHERE phone_number = $1 AND deleted_at is null LIMIT 1`

	updateUser = `
		UPDATE users SET
			employee_id = CASE WHEN $1 THEN $2 ELSE employee_id END,
			phone_number = CASE WHEN $3 THEN $4 ELSE phone_number END,
			first_name = CASE WHEN $5 THEN $6 ELSE first_name END,
			last_name = CASE WHEN $7 THEN $8 ELSE last_name END,
			password_hash = CASE WHEN $9 THEN $10 ELSE password_hash END,
			updated_by = $11,
			updated_at = $12,
			username = CASE WHEN $13 THEN $14 ELSE username END
		WHERE uuid = $15
	`

	findUserById = `
			SELECT
			uuid,
			employee_id,
			username,
			first_name,
			last_name,
			phone_number,
			organization_uuid,
			password_hash,
			created_at,
			is_approved,
			avatar_gradient_start,
			avatar_gradient_end
		FROM users
		WHERE uuid = $1 AND deleted_at is null LIMIT 1
	`

	findByEmployeeId = `
			SELECT
			uuid,
			employee_id,
			username,
			first_name,
			last_name,
			phone_number,
			organization_uuid,
			password_hash,
			created_at,
			is_approved
		FROM users
		WHERE employee_id = $1 AND deleted_at is null LIMIT 1
	`

	listUsers = `
		SELECT
			u.uuid,
			u.employee_id,
			u.username,
			u.first_name,
			u.last_name,
			u.phone_number,
			r.name,
			u.created_at,
			s.name,
			s.uuid
		FROM users u
		JOIN organizations s ON u.organization_uuid = s.uuid
		JOIN roles r ON u.role_uuid = r.uuid
		%s
		%s
		%s
	`

	countUsers = `
		SELECT 
			count(u.uuid) as total_count 
		FROM users u
		JOIN organizations s ON u.organization_uuid = s.uuid
		JOIN roles r ON u.role_uuid = r.uuid
		%s
	`

	deleteUser = `
		UPDATE users SET
			is_deleted = true,
			updated_at = $1,
			updated_by = $2
		WHERE uuid = $3
	`

	insertUserRole = `INSERT INTO user_roles (
		uuid, 
		user_uuid, 
		role_uuid, 
		created_at, 
		created_by,
		updated_at,
		updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING uuid
	`

	deleteUserRole = `DELETE FROM user_roles WHERE uuid = $1`

	deleteUserRoleByUserUUID = `DELETE FROM user_roles WHERE user_uuid = $1`

	findUserRoleById = `
			SELECT 		
			uuid, 
			user_uuid, 
			role_uuid, 
			created_at, 
			created_by
		FROM user_roles
		WHERE uuid = $1 LIMIT 1
	`

	insertPermission = `INSERT INTO permissions (
		uuid, 
		name, 
		action, 
		resource,
		created_at, 
		created_by,
		updated_at,
		updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING uuid
	`

	deletePermission = `DELETE FROM permissions WHERE uuid = $1`

	findPermissionById = `
		SELECT
			uuid, 
			name, 
			action, 
			resource, 
			description,
			created_at, 
			created_by,
			updated_at,
			updated_by
		FROM permissions
		WHERE uuid = $1 LIMIT 1
	`

	findSamePermission = `
		SELECT
			uuid, 
			name, 
			action, 
			resource, 
			description,
			created_at, 
			created_by,
			updated_at,
			updated_by
		FROM permissions
		WHERE lower(name) = lower($1) AND lower(action) = lower($2) AND lower(resource) = lower($3) LIMIT 1
	`

	findSamePermissionExcludeCurrent = `
		SELECT
			uuid, 
			name, 
			action, 
			resource, 
			description,
			created_at, 
			created_by,
			updated_at,
			updated_by
		FROM permissions
		WHERE lower(action) = lower($1) AND lower(resource) = lower($2) AND uuid <> $3 LIMIT 1
	`

	updatePermission = `
		UPDATE permissions SET
			name = CASE WHEN $1 THEN $2 ELSE name END,
			action = CASE WHEN $3 THEN $4 ELSE action END,
			resource = CASE WHEN $5 THEN $6 ELSE resource END,
			description = CASE WHEN $7 THEN $8 ELSE description END,
			updated_by = $9,
			updated_at = $10
		WHERE uuid = $11
	`

	insertRolePermission = `INSERT INTO role_permissions (
		uuid, 
		role_uuid, 
		permission_uuid,
		created_at, 
		created_by,
		updated_at,
		updated_by
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING uuid
	`

	deleteRolePermission = `DELETE FROM role_permissions WHERE uuid = $1`

	findRolePermissionById = `
		SELECT
			uuid,
			role_uuid,
			permission_uuid,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM role_permissions
		WHERE uuid = $1 LIMIT 1
	`

	findRoleByUserId = `
		SELECT 
			r.uuid, 
			r.name,
			r.description
		FROM user_roles u LEFT JOIN roles r on r.uuid = u.role_uuid
		WHERE u.user_uuid = $1 AND deleted_at IS NULL
	`

	findPermissionByRoleUUIDs = `
		SELECT DISTINCT
			p.uuid, p.name, p.action, p.resource, p.description, p.created_at, p.created_by, p.updated_at, p.updated_by
		FROM 
			role_permissions rp
		JOIN 
			permissions p ON rp.permission_uuid = p.uuid
		WHERE 
			rp.role_uuid IN (?)
	`

	findUserRoleByUserUUIDs = `
		SELECT 
			ur.uuid, ur.user_uuid, ur.role_uuid, r.uuid, r.name, r.description, r.created_at, r.created_by, r.updated_at, r.updated_by
		FROM 
			user_roles ur
		LEFT JOIN 
			roles r ON ur.role_uuid = r.uuid
		WHERE 
			ur.user_uuid IN (?)
	`
)
