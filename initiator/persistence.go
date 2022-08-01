package initiator

import (
	"sso/internal/constant/model/db"
	"sso/platform/logger"
)

type Persistence struct {
	// TODO implement
}

func InitPersistence(db *db.Queries, log logger.Logger) Persistence {
	return Persistence{}
}
