package persistencedb

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"sso/internal/constant/model/db"
)

type PersistenceDB struct {
	*db.Queries
	pool *pgxpool.Pool
}

func New(db *db.Queries, pool *pgxpool.Pool) PersistenceDB {
	return PersistenceDB{
		Queries: db,
		pool:    pool,
	}
}
