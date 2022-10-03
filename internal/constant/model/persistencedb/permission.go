package persistencedb

import (
	"context"
	"sso/internal/constant/permissions"
)

const getAllPermissions = "SELECT v0, v1, v2 FROM casbin_rule WHERE p_type = 'p'"

func (db *PersistenceDB) GetAllPermissions(ctx context.Context) ([]permissions.Permission, error) {
	rows, err := db.pool.Query(ctx, getAllPermissions)
	if err != nil {
		return nil, err
	}

	var perms []permissions.Permission
	for rows.Next() {
		var i permissions.Permission
		if err := rows.Scan(&i.ID, &i.Name, &i.Category); err != nil {
			return nil, err
		}
		perms = append(perms, i)
	}

	return perms, nil
}

const getPermissionsOfCategory = "SELECT v0, v1, v2 FROM casbin_rule WHERE p_type = 'p' AND v2 = $1"

func (db *PersistenceDB) GetPermissionsOfCategory(ctx context.Context, category string) ([]permissions.Permission, error) {
	rows, err := db.pool.Query(ctx, getPermissionsOfCategory, category)
	if err != nil {
		return nil, err
	}

	var perms []permissions.Permission
	for rows.Next() {
		var i permissions.Permission
		if err := rows.Scan(&i.ID, &i.Name, &i.Category); err != nil {
			return nil, err
		}
		perms = append(perms, i)
	}

	return perms, nil
}
