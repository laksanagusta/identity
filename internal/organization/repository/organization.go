package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/laksanagusta/identity/internal/entities"
	"github.com/laksanagusta/identity/internal/organization"
	"github.com/laksanagusta/identity/pkg/database"
)

func NewOrganizationRepo(db database.Queryer) organization.Repository {
	return &organizationRepo{
		db: db,
	}
}

type organizationRepo struct {
	db database.Queryer
}

func (r *organizationRepo) WithTransaction(tx database.DBTx) organization.Repository {
	return NewOrganizationRepo(tx)
}

func (r *organizationRepo) Insert(ctx context.Context, organization entities.Organization) (string, error) {
	var returnedUUID string
	err := r.db.GetContext(ctx, &returnedUUID, insertOrganization,
		organization.UUID,
		organization.Name,
		organization.Address,
		organization.Type,
		organization.Path,
		organization.ParentUUID,
		time.Now(),
		organization.CreatedBy,
		time.Now(),
		organization.UpdatedBy,
	)
	if err != nil {
		return "", err
	}
	if returnedUUID != organization.UUID {
		return "", err
	}

	return returnedUUID, nil
}

func (r *organizationRepo) Update(ctx context.Context, organization entities.Organization) error {
	res, err := r.db.ExecContext(ctx,
		updateOrganization,
		organization.Name.IsExists,
		organization.Name,
		organization.Address.IsExists,
		organization.Address,
		organization.Type.IsExists,
		organization.Type,
		time.Now(),
		organization.UpdatedBy,
		organization.UUID,
	)
	if err != nil {
		log.Println("sdawidjwoaidao")
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

func (r *organizationRepo) FindOrganizationByUUID(ctx context.Context, rootUUID string) (*entities.Organization, error) {
	rows, err := r.db.QueryxContext(ctx, findOrganizationById, rootUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodes := make(map[string]*entities.Organization)
	var root *entities.Organization

	for rows.Next() {
		var n entities.Organization
		if err := rows.Scan(
			&n.UUID,
			&n.Name,
			&n.Code,
			&n.Address,
			&n.Type,
			&n.ParentUUID,
			&n.Path,
			&n.Level,
			&n.IsActive,
			&n.CreatedAt,
			&n.CreatedBy,
			&n.UpdatedAt,
			&n.UpdatedBy); err != nil {
			return nil, err
		}

		// Buat copy agar tidak reuse variabel loop
		node := &entities.Organization{
			Name:       n.Name,
			Code:       n.Code,
			Address:    n.Address,
			Type:       n.Type,
			ParentUUID: n.ParentUUID,
			Path:       n.Path,
			Level:      n.Level,
			IsActive:   n.IsActive,
			Children:   []*entities.Organization{},
		}

		node.BaseModel = entities.BaseModel{
			UUID:      n.UUID,
			CreatedAt: n.CreatedAt,
			CreatedBy: n.CreatedBy,
			UpdatedAt: n.UpdatedAt,
			UpdatedBy: n.UpdatedBy,
		}

		nodes[node.UUID] = node

		// Set root
		if node.UUID == rootUUID {
			root = node
		}

		// Hanya tambahkan ke parent jika parent ADA dan BUKAN diri sendiri
		if n.ParentUUID.IsExists && n.ParentUUID.Val != nil {
			parentUUID := *n.ParentUUID.Val
			if parentUUID != node.UUID { // Hindari self-reference
				if parent, ok := nodes[parentUUID]; ok {
					parent.Children = append(parent.Children, node)
				}
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return root, nil
}

func (r *organizationRepo) IndexOrganization(ctx context.Context, params entities.ListOrganizationParams) ([]entities.Organization, *entities.Metadata, error) {
	whereClause := []string{}
	finalArgs := []interface{}{}
	if params.Search.IsExists {
		whereClause = append(whereClause, fmt.Sprintf("lower(s.name) LIKE $%d", len(finalArgs)+1))
		finalArgs = append(finalArgs, "%"+strings.ToLower(params.Search.GetOrDefault())+"%")
	}

	if params.StartTime.IsExists {
		whereClause = append(whereClause, fmt.Sprintf("s.created_by >= $%d", len(finalArgs)+1))
		finalArgs = append(finalArgs, params.StartTime.GetOrDefault())
	}

	if params.EndTime.IsExists {
		whereClause = append(whereClause, fmt.Sprintf("s.created_by <= $%d", len(finalArgs)+1))
		finalArgs = append(finalArgs, params.EndTime.GetOrDefault())
	}

	whereClause = append(whereClause, "s.deleted_at is null")

	whereStr := ""
	if len(whereClause) > 0 {
		whereStr = fmt.Sprintf("WHERE %s", strings.Join(whereClause, " AND "))
	}

	countOrganizationsQuery := fmt.Sprintf(countOrganizations, whereStr)
	var totalCount float64
	row := r.db.QueryRowxContext(ctx, countOrganizationsQuery, finalArgs...)
	err := row.Scan(
		&totalCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &entities.Metadata{}, nil
		}
		return nil, nil, err
	}

	sortStr := ""
	if params.Sort != nil {
		sortStr += fmt.Sprintf("ORDER BY %s %s", params.Sort.FieldName, params.Sort.SortType)
	}

	pagination := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(finalArgs)+1, len(finalArgs)+2)
	finalArgs = append(finalArgs, params.Limit)
	finalArgs = append(finalArgs, params.Offset)

	query := fmt.Sprintf(listOrganization, whereStr, sortStr, pagination)

	var organizations []entities.Organization
	organizationRow, err := r.db.QueryxContext(ctx, query, finalArgs...)
	if err != nil {
		return nil, nil, err
	}

	for organizationRow.Next() {
		var organization entities.Organization
		err := organizationRow.Scan(
			&organization.UUID,
			&organization.Name,
			&organization.Address,
			&organization.CreatedAt,
			&organization.Type,
			&organization.CreatedBy,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, &entities.Metadata{}, nil
			}
			return nil, nil, err
		}

		organizations = append(organizations, organization)
	}

	totalPage := math.Ceil(totalCount / float64(params.Limit))
	if totalPage < 0 {
		totalPage = 0
	}

	return organizations, &entities.Metadata{
		Count:       float64(len(organizations)),
		TotalCount:  totalCount,
		TotalPage:   totalPage,
		CurrentPage: (float64(params.Offset) / float64(params.Limit)) + 1,
	}, nil
}

func (r *organizationRepo) Delete(ctx context.Context, uuid string, username string) error {
	res, err := r.db.ExecContext(ctx,
		deleteOrganization,
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

func (r *organizationRepo) FindOrganizationByUUIDs(ctx context.Context, uuids []string) ([]*entities.Organization, error) {
	if len(uuids) == 0 {
		return []*entities.Organization{}, nil
	}

	query, args, err := sqlx.In(findOrganizationUUIDs, uuids)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	rows, err := r.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []*entities.Organization
	for rows.Next() {
		var organization entities.Organization
		err := rows.Scan(
			&organization.UUID,
			&organization.Name,
			&organization.Code,
			&organization.Address,
			&organization.Type,
			&organization.CreatedAt,
			&organization.CreatedBy,
			&organization.UpdatedAt,
			&organization.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		organizations = append(organizations, &organization)
	}

	return organizations, nil
}
