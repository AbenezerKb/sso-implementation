package initiator

import (
	"sso/internal/constant/model/db"
	"sso/internal/storage"
	"sso/internal/storage/oauth"
	"sso/platform/logger"
)

type Persistence struct {
	// TODO implement
	OAuthPersistence storage.OAuthPersistence
}

func InitPersistence(db *db.Queries, log logger.Logger) Persistence {
	return Persistence{
		OAuthPersistence: oauth.InitOAuth(log, db),
	}
}
