package repository

var (
	findRoleById = `SELECT uuid, name FROM roles WHERE uuid = $1 AND deleted_at IS NULL`

	findRole = `SELECT uuid, name, description FROM roles WHERE deleted_at IS NULL`

	insertRole = `INSERT INTO roles (
		uuid, 
		name, 
		description, 
		created_at, 
		created_by,
		updated_at,
		updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING uuid
	`

	deleteRole = `DELETE FROM roles WHERE uuid = $1`

	findRoleByName = `SELECT 		
		uuid, 
		name
	FROM roles
	WHERE LOWER(name) = LOWER($1) AND deleted_at is null LIMIT 1`
)
