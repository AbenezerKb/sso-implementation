package db

import (
	"context"

	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
)

func (q *Queries) GetAllClients(ctx context.Context, pgnFlt db_pgnflt.FilterParams) ([]Client, int, error) {
	_, sql := db_pgnflt.GetFilterSQL(pgnFlt)
	rows, err := q.db.Query(ctx, db_pgnflt.GetSelectColumnsQuery([]string{
		"id",
		"name",
		"client_type",
		"redirect_uris",
		"scopes",
		"secret",
		"logo_url",
		"status",
		"created_at",
		"first_party",
	}, "clients", sql))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var clients []Client
	var totalCount int
	for rows.Next() {
		var i Client
		if err := rows.Scan(&i.ID, &i.Name,
			&i.ClientType,
			&i.RedirectUris,
			&i.Scopes,
			&i.Secret,
			&i.LogoUrl,
			&i.Status,
			&i.CreatedAt,
			&i.FirstParty,
			&totalCount); err != nil {
			return nil, 0, err
		}
		clients = append(clients, i)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return clients, totalCount, nil
}
