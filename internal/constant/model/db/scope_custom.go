package db

import (
	"context"

	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
)

func (q *Queries) GetAllScopes(ctx context.Context, pgnFlt db_pgnflt.FilterParams) ([]Scope, int, error) {
	_, sql := db_pgnflt.GetFilterSQL(pgnFlt)
	rows, err := q.db.Query(ctx, db_pgnflt.GetSelectColumnsQuery([]string{
		"id",
		"name",
		"description",
		"resource_server_id",
		"resource_server_name",
		"status",
		"created_at",
	}, "scopes", sql))
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
