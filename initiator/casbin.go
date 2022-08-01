package initiator

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/pckhoi/casbin-pgx-adapter"
	"sso/platform/logger"
)

func InitEnforcer(path, conn string, log logger.Logger) *casbin.Enforcer {
	adapter, err := pgxadapter.NewAdapter(conn) // FIXME: check casbin table
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to create adapter: %v", err))
	}

	enforcer, err := casbin.NewEnforcer(path, adapter)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to create enforcer: %v", err))
	}

	return enforcer
}
