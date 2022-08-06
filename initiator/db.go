package initiator

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"sso/internal/constant/model/db"
	"sso/platform/logger"
)

func InitDB(url string, log logger.Logger) *db.Queries {
	config, err := pgx.ParseConfig(url)
	config.Logger = log
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}
	return db.New(conn)
}
