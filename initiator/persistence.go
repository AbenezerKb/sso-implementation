package initiator

import (
	"sso/internal/constant/model/db"
	"sso/internal/storage"
	"sso/internal/storage/persistence/client"
	"sso/internal/storage/persistence/oauth"
	"sso/internal/storage/persistence/oauth2"
	"sso/platform/logger"
)

type Persistence struct {
	// TODO implement
	OAuthPersistence  storage.OAuthPersistence
	OAuth2Persistence storage.OAuth2Persistence
	ClientPersistence storage.ClientPersistence
}

func InitPersistence(db *db.Queries, log logger.Logger) Persistence {
	return Persistence{
		OAuthPersistence:  oauth.InitOAuth(log, db),
		OAuth2Persistence: oauth2.InitOAuth2(log, db),
		ClientPersistence: client.InitClient(log, db),
	}
}
