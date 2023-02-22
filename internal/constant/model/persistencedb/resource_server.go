package persistencedb

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"

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

func (p *PersistenceDB) GetAllResourceServers(ctx context.Context, pgnFlt db_pgnflt.FilterParams) ([]dto.ResourceServer, int, error) {
	_, sqlStr := db_pgnflt.GetFilterSQL(pgnFlt)
	rows, err := p.pool.Query(ctx, db_pgnflt.GetSelectColumnsQueryWithJoins([]string{

		"rs.id AS resource_server_id",
		"rs.name AS resource_server_name",
		"rs.created_at",
		"rs.updated_at",
		"sc.id AS scope_id",
		"sc.name AS scope_name",
		"sc.description",
		"sc.status",
	},
		db_pgnflt.Table{Name: "resource_servers", Alias: "rs"}, []db_pgnflt.JOIN{
			{
				JoinType: "LEFT JOIN",
				Table: db_pgnflt.Table{
					Name:  "scopes",
					Alias: "sc",
				},
				On: "sc.resource_server_name = rs.name",
			},
		}, sqlStr))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// maps are a better way to search than slices
	resourceServers := map[uuid.UUID]dto.ResourceServer{}
	var totalCount, reducer int
	for rows.Next() {
		var i db2.ResourceServer
		var s struct {
			ID          uuid.UUID
			Name        sql.NullString
			Description sql.NullString
			Status      sql.NullString
		}
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedAt,
			&i.UpdatedAt,
			&s.ID,
			&s.Name,
			&s.Description,
			&s.Status,
			&totalCount); err != nil {
			return nil, 0, err
		}
		if v, ok := resourceServers[i.ID]; ok {
			v.Scopes = append(v.Scopes, dto.Scope{
				Name:        s.Name.String,
				Description: s.Description.String,
			})
			resourceServers[i.ID] = v
			reducer++
		} else {
			rs := dto.ResourceServer{
				ID:        i.ID,
				Name:      i.Name,
				CreatedAt: i.CreatedAt,
				UpdatedAt: i.UpdatedAt,
			}
			if s.Name.Valid { // if scope was found
				rs.Scopes = []dto.Scope{
					{
						Name:        s.Name.String,
						Description: s.Description.String,
					},
				}
			}
			resourceServers[i.ID] = rs
		}
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	servers := make([]dto.ResourceServer, 0, len(resourceServers))
	for _, v := range resourceServers {
		servers = append(servers, v)
	}
	return servers, totalCount - reducer, nil
}
