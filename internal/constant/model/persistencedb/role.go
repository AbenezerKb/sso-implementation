package persistencedb

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"sso/internal/constant/errors/sqlcerr"
	db2 "sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
)

const getRoleForUser = "SELECT v1 FROM casbin_rule WHERE v0 = $1"

func (db *PersistenceDB) GetRoleForUser(ctx context.Context, userID uuid.UUID) (string, error) {
	row := db.pool.QueryRow(ctx, getRoleForUser, userID)
	var role string
	err := row.Scan(&role)
	if err != nil {
		return "", err
	}

	return role, nil
}

const createRole = "INSERT INTO casbin_rule (p_type, v0, v1, v2) values ('g', $1, $2, 'role') RETURNING v1"

func (db *PersistenceDB) CreateRoleTX(ctx context.Context, roleName string, perms []string) (dto.Role, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return dto.Role{}, err
	}
	defer func(ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(ctx)

	query := db.Queries.WithTx(tx)
	var dbPerms []string
	for i := 0; i < len(perms); i++ {
		row := tx.QueryRow(ctx, createRole, roleName, perms[i])
		var perm string
		if err := row.Scan(&perm); err != nil {
			return dto.Role{}, err
		}
		dbPerms = append(dbPerms, perm)
	}

	dbRole, err := query.AddRole(ctx, roleName)
	if err != nil {
		return dto.Role{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return dto.Role{}, err
	}

	return dto.Role{
		Name:        dbRole.Name,
		Permissions: dbPerms,
	}, nil
}

const checkIfPermissionExists = "SELECT v0 FROM casbin_rule WHERE v0 = $1"

func (db *PersistenceDB) CheckIfPermissionExists(ctx context.Context, permission string) (bool, error) {
	row := db.pool.QueryRow(ctx, checkIfPermissionExists, permission)
	var i string
	err := row.Scan(&i)
	if err != nil {
		if sqlcerr.Is(err, sqlcerr.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

const getAllRoles = `
SELECT r.name,
       r.status,
       r.created_at,
       r.updated_at,
       (SELECT string_to_array(string_agg(v1, ','), ',')
        FROM casbin_rule
        WHERE v0 = r.name) AS permissions,
       count(*) over()
FROM roles r`

func (db *PersistenceDB) GetAllRoles(ctx context.Context, pgnFlt string) ([]dto.Role, int, error) {
	rows, err := db.pool.Query(ctx, fmt.Sprintf("%s %s", getAllRoles, pgnFlt))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// maps are a better way to search than slices
	var roles []dto.Role
	var totalCount int
	for rows.Next() {
		var i db2.Role
		var p []string
		if err := rows.Scan(
			&i.Name,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&p,
			&totalCount); err != nil {
			return nil, 0, err
		}
		roles = append(roles, dto.Role{
			Name:        i.Name,
			Status:      i.Status.String,
			CreatedAt:   i.CreatedAt,
			UpdatedAt:   i.UpdatedAt,
			Permissions: p,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return roles, totalCount, nil
}
