package persistencedb

import (
	"context"
	"sso/internal/constant/model/db"
	"sso/platform/utils"
)

func (q *PersistenceDB) GetAllIdentityProviders(ctx context.Context, pgnFlt string) ([]db.IdentityProvider, int, error) {
	rows, err := q.pool.Query(ctx, utils.ComposeFullFilterSQL(ctx, "identity_providers", pgnFlt))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var idPs []db.IdentityProvider
	var totalCount int
	for rows.Next() {
		var i db.IdentityProvider
		if err := rows.Scan(&i.ID, &i.Name,
			&i.AuthorizationUri,
			&i.RedirectUri,
			&i.ClientID,
			&i.ClientSecret,
			&i.TokenEndpointUrl,
			&i.UserInfoEndpointUrl,
			&i.LogoUrl,
			&i.Status,
			&i.UpdatedAt,
			&i.CreatedAt,
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