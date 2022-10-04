package persistencedb

import (
	"context"
	"github.com/google/uuid"
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
