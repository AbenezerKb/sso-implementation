package initiator

import (
	"sso/internal/constant/model/persistencedb"
	"sso/internal/storage"
	"sso/internal/storage/persistence/client"
	"sso/internal/storage/persistence/oauth"
	"sso/internal/storage/persistence/oauth2"
	"sso/internal/storage/persistence/profile"
	"sso/internal/storage/persistence/scope"
	"sso/internal/storage/persistence/user"
	"sso/platform/logger"
)

type Persistence struct {
	// TODO implement
	OAuthPersistence    storage.OAuthPersistence
	OAuth2Persistence   storage.OAuth2Persistence
	ClientPersistence   storage.ClientPersistence
	ScopePersistence    storage.ScopePersistence
	UserPersistence     storage.UserPersistence
	ProfilePersistence  storage.ProfilePersistence
	MiniRidePersistence storage.MiniRidePersistence
}

func InitPersistence(db persistencedb.PersistenceDB, log logger.Logger) Persistence {
	return Persistence{
		OAuthPersistence:   oauth.InitOAuth(log.Named("oauth-persistence"), db.Queries),
		ClientPersistence:  client.InitClient(log.Named("client-persistence"), db.Queries),
		OAuth2Persistence:  oauth2.InitOAuth2(log.Named("oauth2-persistence"), db.Queries),
		ScopePersistence:   scope.InitScopePersistence(log.Named("scope-persistence"), db.Queries),
		UserPersistence:    user.InitUserPersistence(log.Named("user-persistence"), db.Queries),
		ProfilePersistence: profile.InitProfilePersistence(log.Named("profile-persistence"), db.Queries),
	}
}
