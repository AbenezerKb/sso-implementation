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
	checkPeriod := viper.GetDuration("database.health_check_period")
	if checkPeriod == 0 {
		checkPeriod = 5 * time.Minute
	}
	config.HealthCheckPeriod = checkPeriod

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}
	return conn
}
