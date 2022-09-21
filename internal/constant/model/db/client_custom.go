package db

import (
	"context"
	"sso/platform/utils"
)

func (q *Queries) GetAllClients(ctx context.Context, pgnFlt string) ([]Client, int, error) {
	rows, err := q.db.Query(ctx, utils.ComposeFullFilterSQL(ctx, "clients", pgnFlt))
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
