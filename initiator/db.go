package initiator

import (
	"context"
	"fmt"
	"time"

	"sso/platform/logger"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

func InitDB(url string, log logger.Logger) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}
	config.ConnConfig.Logger = log.Named("pgx")

	idleConnTimeout := viper.GetDuration("database.idle_conn_timeout")
	if idleConnTimeout == 0 {
		idleConnTimeout = 4 * time.Minute
	}
	config.MaxConnIdleTime = idleConnTimeout

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}
	return conn
}
