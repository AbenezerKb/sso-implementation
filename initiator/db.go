package initiator

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"sso/internal/constant/model/db"
	"sso/platform/logger"
)

func InitDB(url string, log logger.Logger) *db.Queries {
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("Failed to connect to database: %v", err))
	}
	return db.New(conn)
}
