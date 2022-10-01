package persistencedb

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4"
	db2 "sso/internal/constant/model/db"
	"sso/internal/constant/model/dto"
)

func (db *PersistenceDB) CreateResourceServerWithTX(ctx context.Context, server dto.ResourceServer) (dto.ResourceServer, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return dto.ResourceServer{}, err
	}
	defer func(tx pgx.Tx) {
		_ = tx.Rollback(ctx) // FIXME: a better way to handle this
	}(tx)

	query := db.Queries.WithTx(tx)
	createdServer, err := query.CreateResourceServer(ctx, server.Name)
	if err != nil {
		return dto.ResourceServer{}, err
	}

	var scopes []dto.Scope
	for i := 0; i < len(server.Scopes); i++ {
		scope, err := query.CreateScope(ctx, db2.CreateScopeParams{
			Name:        server.Scopes[i].Name,
			Description: server.Scopes[i].Description,
			ResourceServerName: sql.NullString{
				String: server.Name,
				Valid:  true,
			},
		})
		if err != nil {
			return dto.ResourceServer{}, err
		}
		scopes = append(scopes, dto.Scope{
			Name:               scope.Name,
			Description:        scope.Description,
			ResourceServerName: scope.ResourceServerName.String,
		})
	}

	if err := tx.Commit(ctx); err != nil {
		return dto.ResourceServer{}, err
	}
	return dto.ResourceServer{
		ID:        createdServer.ID,
		Name:      createdServer.Name,
		CreatedAt: createdServer.CreatedAt,
		UpdatedAt: createdServer.UpdatedAt,
		Scopes:    scopes,
	}, nil
}
