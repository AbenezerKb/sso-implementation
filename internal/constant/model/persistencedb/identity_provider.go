package persistencedb

import (
	"context"

	"sso/internal/constant/model/db"

	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
)

func (q *PersistenceDB) GetAllIdentityProviders(ctx context.Context, pgnFlt db_pgnflt.FilterParams) ([]db.IdentityProvider, int, error) {
	_, sql := db_pgnflt.GetFilterSQL(pgnFlt)
	rows, err := q.pool.Query(ctx, db_pgnflt.GetSelectColumnsQuery([]string{
		"id",
		"name",
		"logo_url",
		"client_id",
		"client_secret",
		"redirect_uri",
		"authorization_uri",
		"token_endpoint_url",
		"user_info_endpoint_url",
		"status",
		"created_at",
		"updated_at",
	}, "identity_providers", sql))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var idPs []db.IdentityProvider
	var totalCount int
	for rows.Next() {
		var i db.IdentityProvider
		if err := rows.Scan(&i.ID, &i.Name,
			&i.LogoUrl,
			&i.ClientID,
			&i.ClientSecret,
			&i.RedirectUri,
			&i.AuthorizationUri,
			&i.TokenEndpointUrl,
			&i.UserInfoEndpointUrl,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&totalCount); err != nil {
			return nil, 0, err
		}
		idPs = append(idPs, i)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return idPs, totalCount, nil
}
