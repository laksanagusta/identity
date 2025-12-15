package repository

var (
	insertOrganization = `INSERT INTO organizations (
		uuid,
		name,
		code,
		address,
		type,
		path,
		parent_uuid,
		created_at,
		created_by,
		updated_at,
		updated_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING UUID`

	updateOrganization = `
		UPDATE organizations SET
			name = CASE WHEN $1 THEN $2 ELSE name END,
			address = CASE WHEN $3 THEN $4 ELSE address END,
			type = CASE WHEN $5 THEN $6 ELSE type END,
			updated_at = $7,
			updated_by = $8
		WHERE uuid = $9
	`

	findOrganizationById = `
		WITH target_org AS (
			SELECT uuid, path FROM organizations WHERE uuid = $1 AND deleted_at IS NULL
		)
		SELECT
			o.uuid,
			o.name,
			o.code,
			o.address,
			o.type,
			o.parent_uuid,
			o.path,
			o.level,
			o.is_active,
			o.created_at,
			o.created_by,
			o.updated_at,
			o.updated_by
		FROM organizations o, target_org t
		WHERE o.deleted_at IS NULL 
			AND (o.uuid = t.uuid OR o.path LIKE t.path || '.%')
		ORDER BY o.path
	`

	listOrganization = `
		SELECT 		
			s.uuid, 
			s.name, 
			s.address,
			s.created_at,
			s.type,
			s.created_by
		FROM organizations s
		%s
		%s
		%s
	`

	countOrganizations = `
		SELECT 
			count(s.uuid) as total_count 
		FROM organizations s
		%s
	`

	deleteOrganization = `
		UPDATE organizations SET is_deleted = true, updated_at = $1, updated_by = $2 WHERE uuid = $3
	`

	findOrganizationUUIDs = `
		SELECT 
			o.uuid, 
			o.name, 
			o.code, 
			o.address, 
			o.type, 
			o.created_at, 
			o.created_by, 
			o.updated_at, 
			o.updated_by
		FROM 
			organizations o
		WHERE 
			o.uuid IN (?)
	`
)
