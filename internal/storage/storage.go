package storage

import (
	"context"

	"sso/internal/constant/model"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"

	"github.com/google/uuid"
	db_pgnflt "gitlab.com/2ftimeplc/2fbackend/repo/db-pgnflt"
)

type OAuthPersistence interface {
	Register(ctx context.Context, user dto.User) (*dto.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*dto.User, error)
	GetUserStatus(ctx context.Context, Id uuid.UUID) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*dto.User, error)
	UserByPhoneExists(ctx context.Context, phone string) (bool, error)
	UserByEmailExists(ctx context.Context, email string) (bool, error)
	GetUserByPhoneOrEmail(ctx context.Context, query string) (*dto.User, error)
	GetUserByID(ctx context.Context, Id uuid.UUID) (*dto.User, error)
	RemoveInternalRefreshToken(ctx context.Context, refreshToken string) error
	SaveInternalRefreshToken(ctx context.Context, rf dto.InternalRefreshToken) error
	GetInternalRefreshToken(ctx context.Context, refreshtoken string) (*dto.InternalRefreshToken, error)
	UpdateInternalRefreshToken(ctx context.Context, oldToken, newToken string) (*dto.InternalRefreshToken, error)
	GetInternalRefreshTokensByUserID(ctx context.Context, userID uuid.UUID) ([]dto.InternalRefreshToken, error)
	GetUserPassword(ctx context.Context, Id uuid.UUID) (string, error)
	GetAllIdentityProviders(ctx context.Context) ([]dto.IdentityProvider, error)
	ChangeUserPassword(ctx context.Context, phone, newPassword string) error
}

type OTPCache interface {
	SetOTP(ctx context.Context, phone string, otp string) error
	GetOTP(ctx context.Context, phone string) (string, error)
	GetDelOTP(ctx context.Context, phone string) (string, error)
	DeleteOTP(ctx context.Context, phone ...string) error
	VerifyOTP(ctx context.Context, phone string, otp string) error
}

type SessionCache interface {
	SaveSession(ctx context.Context, session dto.Session) error
	GetSession(ctx context.Context, sessionID string) (dto.Session, error)
}

type OAuth2Persistence interface {
	GetNamedScopes(ctx context.Context, scopes ...string) ([]dto.Scope, error)
	AuthHistoryExists(ctx context.Context, code string) (bool, error)
	PersistRefreshToken(ctx context.Context, param dto.RefreshToken) (*dto.RefreshToken, error)
	RemoveRefreshTokenCode(ctx context.Context, code string) error
	RemoveRefreshToken(ctx context.Context, refresh_token string) error
	AddAuthHistory(ctx context.Context, param dto.AuthHistory) (*dto.AuthHistory, error)
	CheckIfUserGrantedClient(ctx context.Context, userID uuid.UUID, clientID uuid.UUID) (bool, dto.RefreshToken, error)
	GetRefreshToken(ctx context.Context, token string) (*dto.RefreshToken, error)
	GetRefreshTokenOfClientByUserID(ctx context.Context, userID, clientID uuid.UUID) (*dto.RefreshToken, error)
	GetAuthorizedClients(ctx context.Context, userID uuid.UUID) ([]dto.AuthorizedClientsResponse, error)
	GetOpenIDAuthorizedClients(ctx context.Context, userID uuid.UUID) ([]dto.AuthorizedClientsResponse, error)
	UserInfo(ctx context.Context, userID uuid.UUID) (*dto.UserInfo, error)
	UpdateRefreshToken(ctx context.Context, newRefreshToken, oldRefreshToken string) (*dto.RefreshToken, error)
}

type ConsentCache interface {
	SaveConsent(ctx context.Context, consent dto.Consent) error
	GetConsent(ctx context.Context, consentID string) (dto.Consent, error)
	DeleteConsent(ctx context.Context, consentID string) error
	ChangeStatus(ctx context.Context, status bool, consent dto.Consent) (dto.Consent, error)
}
type ClientPersistence interface {
	Create(ctx context.Context, client dto.Client) (*dto.Client, error)
	GetClientByID(ctx context.Context, id uuid.UUID) (*dto.Client, error)
	DeleteClientByID(ctx context.Context, id uuid.UUID) error
	GetAllClients(ctx context.Context, filters db_pgnflt.FilterParams) ([]dto.Client, *model.MetaData, error)
	UpdateClientStatus(ctx context.Context, updateClientStatusParam dto.UpdateClientStatus, clientID uuid.UUID) error
	UpdateClient(ctx context.Context, client dto.Client) error
}

type AuthCodeCache interface {
	SaveAuthCode(ctx context.Context, authCode dto.AuthCode) error
	GetAuthCode(ctx context.Context, code string) (dto.AuthCode, error)
	DeleteAuthCode(ctx context.Context, code string) error
}

type ResetCodeCache interface {
	SaveResetCode(ctx context.Context, email, code string) error
	GetResetCode(ctx context.Context, email string) (string, error)
	DeleteResetCode(ctx context.Context, email string) error
}

type ScopePersistence interface {
	CreateScope(ctx context.Context, scope dto.Scope) (dto.Scope, error)
	GetScope(ctx context.Context, scope string) (dto.Scope, error)
	GetListedScopes(ctx context.Context, scopes ...string) ([]dto.Scope, error)
	GetScopeNameOnly(ctx context.Context, scopes ...string) (string, error)
	GetAllScopes(ctx context.Context, filters db_pgnflt.FilterParams) ([]dto.Scope, *model.MetaData, error)
	DeleteScopeByName(ctx context.Context, name string) error
	UpdateScope(ctx context.Context, scopeUpdateParam dto.Scope) error
}

type UserPersistence interface {
	GetAllUsers(ctx context.Context, filters db_pgnflt.FilterParams) ([]dto.User, *model.MetaData, error)
	UpdateUserStatus(ctx context.Context, updateUserStatusParam dto.UpdateUserStatus, userID uuid.UUID) error
	UpdateUserRole(ctx context.Context, userID uuid.UUID, roleName string) error
	RevokeUserRole(ctx context.Context, userID uuid.UUID) error
	GetUserByID(ctx context.Context, Id uuid.UUID) (*dto.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*dto.User, error)
	GetUsersByPhone(ctx context.Context, phones []string) ([]dto.User, error)
	GetUsersByID(ctx context.Context, ids []string) ([]dto.User, error)
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) (*dto.User, error)
	DeleteUser(ctx context.Context, UserID uuid.UUID) error
}

type ProfilePersistence interface {
	UpdateProfile(ctx context.Context, userParam dto.User) (*dto.User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*dto.User, error)
	UpdateProfilePicture(ctx context.Context, finalImageName string, userID uuid.UUID) error
	ChangePhone(ctx context.Context, changePhoneParam dto.ChangePhoneParam, userID uuid.UUID) error
	ChangePassword(ctx context.Context, changePasswordParam dto.ChangePasswordParam, userID uuid.UUID) error
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type ResourceServerPersistence interface {
	CreateResourceServer(ctx context.Context, server dto.ResourceServer) (dto.ResourceServer, error)
	GetResourceServerByName(ctx context.Context, name string) (dto.ResourceServer, error)
	GetAllResourceServers(ctx context.Context, filters db_pgnflt.FilterParams) ([]dto.ResourceServer, *model.MetaData, error)
	GetResourceServerByID(ctx context.Context, rsID uuid.UUID) (*dto.ResourceServer, error)
}

type MiniRidePersistence interface {
	UpdateUser(ctx context.Context, updateUserParam *request_models.Driver) error
	CreateUser(ctx context.Context, createUserParam *request_models.Driver) (*dto.User, error)
	SwapPhones(ctx context.Context, newPhone, oldPhone string) error
	CheckPhone(ctx context.Context, phone string) (*dto.MiniRideResponse, error)
}

type RolePersistence interface {
	GetAllPermissions(ctx context.Context, category string) ([]dto.Permission, error)
	GetRoleStatus(ctx context.Context, roleName string) (string, error)
	GetRoleForUser(ctx context.Context, userID uuid.UUID) (string, error)
	GetRoleStatusForUser(ctx context.Context, userID uuid.UUID) (string, error)
	CreateRole(ctx context.Context, role dto.Role) (dto.Role, error)
	CheckIfPermissionExists(ctx context.Context, permission string) (bool, error)
	GetAllRoles(ctx context.Context, filters db_pgnflt.FilterParams) ([]dto.Role, *model.MetaData, error)
	GetRoleByName(ctx context.Context, roleName string) (dto.Role, error)
	UpdateRoleStatus(ctx context.Context, updateStatusParam dto.UpdateRoleStatus, roleName string) error
	DeleteRole(ctx context.Context, roleName string) error
	UpdateRole(ctx context.Context, role dto.UpdateRole) (dto.Role, error)
}

type IdentityProviderPersistence interface {
	CreateIdentityProvider(ctx context.Context, provider dto.IdentityProvider) (dto.IdentityProvider, error)
	GetIdentityProvider(ctx context.Context, ipID uuid.UUID) (dto.IdentityProvider, error)
	SaveIPAccessToken(ctx context.Context, ipAccessToken dto.IPAccessToken) (dto.IPAccessToken, error)
	GetIPAccessTokenBySubAndIP(ctx context.Context, subID string, ipID uuid.UUID) (dto.IPAccessToken, error)
	UpdateIpAccessToken(ctx context.Context, ipAccessToken dto.IPAccessToken) (dto.IPAccessToken, error)
	UpdateIdentityProvider(ctx context.Context, idPParam dto.IdentityProvider) error
	DeleteIdentityProvider(ctx context.Context, idPID uuid.UUID) error
	GetAllIdentityProviders(ctx context.Context, filters db_pgnflt.FilterParams) ([]dto.IdentityProvider, *model.MetaData, error)
}
