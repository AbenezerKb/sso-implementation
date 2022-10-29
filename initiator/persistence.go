package initiator

import (
	"sso/internal/constant/model/persistencedb"
	"sso/internal/storage"
	"sso/internal/storage/persistence/client"
	identity_provider "sso/internal/storage/persistence/identity-provider"
	"sso/internal/storage/persistence/mini_ride"
	"sso/internal/storage/persistence/oauth"
	"sso/internal/storage/persistence/oauth2"
	"sso/internal/storage/persistence/profile"
	resource_server "sso/internal/storage/persistence/resource-server"
	"sso/internal/storage/persistence/role"
	"sso/internal/storage/persistence/scope"
	"sso/internal/storage/persistence/user"
	"sso/platform/logger"
)

type Persistence struct {
	// TODO implement
	OAuthPersistence            storage.OAuthPersistence
	OAuth2Persistence           storage.OAuth2Persistence
	ClientPersistence           storage.ClientPersistence
	ScopePersistence            storage.ScopePersistence
	UserPersistence             storage.UserPersistence
	ProfilePersistence          storage.ProfilePersistence
	ResourceServerPersistence   storage.ResourceServerPersistence
	MiniRidePersistence         storage.MiniRidePersistence
	RolePersistence             storage.RolePersistence
	IdentityProviderPersistence storage.IdentityProviderPersistence
}

func InitPersistence(db persistencedb.PersistenceDB, log logger.Logger) Persistence {
	return Persistence{
		OAuthPersistence:            oauth.InitOAuth(log.Named("oauth-persistence"), db.Queries),
		ClientPersistence:           client.InitClient(log.Named("client-persistence"), db.Queries),
		OAuth2Persistence:           oauth2.InitOAuth2(log.Named("oauth2-persistence"), db.Queries),
		ScopePersistence:            scope.InitScopePersistence(log.Named("scope-persistence"), db.Queries),
		UserPersistence:             user.InitUserPersistence(log.Named("user-persistence"), &db),
		ProfilePersistence:          profile.InitProfilePersistence(log.Named("profile-persistence"), &db),
		ResourceServerPersistence:   resource_server.InitResourceServerPersistence(log.Named("resource-server-persistence"), &db),
		MiniRidePersistence:         mini_ride.InitMiniRidePersistence(log.Named("mini-ride-persistence"), &db),
		RolePersistence:             role.InitRolePersistence(log.Named("role-persistence"), &db),
		IdentityProviderPersistence: identity_provider.InitIdentityProviderPersistence(log.Named("identity-provider-persistence"), &db),
	}
}
