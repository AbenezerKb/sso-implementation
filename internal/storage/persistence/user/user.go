package user

import (
	"sso/internal/constant/model/db"
	"sso/internal/storage"
	"sso/platform/logger"
)

type userPersistence struct {
	logger logger.Logger
	db     *db.Queries
}

func InitUserPersistence(logger logger.Logger, db *db.Queries) storage.UserPersistence {
	return &userPersistence{
		logger: logger,
		db:     db,
	}
}
