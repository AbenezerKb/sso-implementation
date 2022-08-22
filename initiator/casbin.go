package initiator

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"sso/platform/logger"
	"sso/platform/pgxadapter"
)

func InitEnforcer(path string, conn *pgxpool.Pool, log logger.Logger) *casbin.Enforcer {
	pgxConn, err := conn.Acquire(context.Background())
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}
	adapter, err := pgxadapter.NewAdapterWithDB(pgxConn.Conn())
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to create adapter: %v", err))
	}

	enforcer, err := casbin.NewEnforcer(path, adapter)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to create enforcer: %v", err))
	}

	return enforcer
}
