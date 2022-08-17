package initiator

import (
	"sso/internal/constant/model/db"
	"sso/internal/storage"
	"sso/internal/storage/persistence/client"
	"sso/internal/storage/persistence/oauth"
	"sso/platform/logger"
)

type Persistence struct {
	OAuthPersistence  storage.OAuthPersistence
	ClientPersistence storage.ClientPersistence
}

func InitPersistence(db *db.Queries, log logger.Logger) Persistence {
	return Persistence{
		OAuthPersistence:  oauth.InitOAuth(log.Named("oauth-persistence"), db),
		ClientPersistence: client.InitClient(log.Named("client-persistence"), db),
	}
}
