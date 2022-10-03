package db

import (
	"context"
	"sso/platform/utils"
)

func (q *Queries) GetAllScopes(ctx context.Context, pgnFlt string) ([]Scope, int, error) {
	rows, err := q.db.Query(ctx, utils.ComposeFullFilterSQL(ctx, "scopes", pgnFlt))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var scopes []Scope
	var totalCount int
	for rows.Next() {
		var i Scope
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.ResourceServerID,
			&i.ResourceServerName,
			&i.Status,
			&i.CreatedAt,
			&totalCount); err != nil {
			return nil, 0, err
		}
		scopes = append(scopes, i)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return scopes, totalCount, nil
}
