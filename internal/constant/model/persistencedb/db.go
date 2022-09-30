package persistencedb

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"sso/internal/constant/model/db"
)

type PersistenceDB struct {
	*db.Queries
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) PersistenceDB {
	return PersistenceDB{
		Queries: db.New(pool),
		pool:    pool,
	}
}
