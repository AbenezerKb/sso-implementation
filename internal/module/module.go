package module

import (
	"context"
	"mime/multipart"
	"sso/internal/constant/model"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sync"

	"github.com/joomcode/errorx"

	"github.com/google/uuid"
)

type OAuthModule interface {
	Register(ctx context.Context, user dto.RegisterUser) (*dto.User, error)
	Login(ctx context.Context, login dto.LoginCredential) (*dto.TokenResponse, error)
	ComparePassword(hashedPwd, plainPassword string) bool
	RequestOTP(ctx context.Context, phone string, rqType string) error
	GetUserStatus(ctx context.Context, Id string) (string, error)
	Logout(ctx context.Context, param dto.InternalRefreshTokenRequestBody) error
	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error)
}

type OAuth2Module interface {
	Authorize(ctx context.Context, authRequestParma dto.AuthorizationRequestParam, requestOrigin string, bindError *errorx.Error) string
	GetConsentByID(ctx context.Context, consentID string) (dto.ConsentResponse, error)
	ApproveConsent(ctx context.Context, consentID string, userID uuid.UUID, opbs string, bindError *errorx.Error) string
	RejectConsent(ctx context.Context, consentID, failureReason string, bindError *errorx.Error) string
	Token(ctx context.Context, client dto.Client, param dto.AccessTokenRequest) (*dto.TokenResponse, error)
	Logout(ctx context.Context, logoutReqParam dto.LogoutRequest, bindError *errorx.Error) string
	RevokeClient(ctx context.Context, clientBody request_models.RevokeClientBody) error
	GetAuthorizedClients(ctx context.Context) ([]dto.AuthorizedClientsResponse, error)
	GetOpenIDAuthorizedClients(ctx context.Context) ([]dto.AuthorizedClientsResponse, error)
}
type UserModule interface {
	Create(ctx context.Context, user dto.CreateUser) (*dto.User, error)
	GetUserByID(ctx context.Context, id string) (*dto.User, error)
	GetAllUsers(ctx context.Context, filtersQuery request_models.PgnFltQueryParams) ([]dto.User, *model.MetaData, error)
	UpdateUserStatus(ctx context.Context, updateUserStatusParam dto.UpdateUserStatus, userID string) error
}

type ClientModule interface {
	Create(ctx context.Context, client dto.Client) (*dto.Client, error)
	GetClientByID(ctx context.Context, id string) (*dto.Client, error)
	DeleteClientByID(ctx context.Context, id string) error
	GetAllClients(ctx context.Context, filtersQuery request_models.PgnFltQueryParams) ([]dto.Client, *model.MetaData, error)
}

type ScopeModule interface {
	GetScope(ctx context.Context, scope string) (dto.Scope, error)
	CreateScope(ctx context.Context, scope dto.Scope) (dto.Scope, error)
	GetAllScopes(ctx context.Context, filtersQuery request_models.PgnFltQueryParams) ([]dto.Scope, *model.MetaData, error)
}

type ProfileModule interface {
	UpdateProfile(ctx context.Context, userParam dto.User) (*dto.User, error)
	GetProfile(ctx context.Context) (*dto.User, error)
	UpdateProfilePicture(ctx context.Context, imageFile *multipart.FileHeader) error
}

type ResourceServerModule interface {
	CreateResourceServer(ctx context.Context, server dto.ResourceServer) (dto.ResourceServer, error)
	GetAllResourceServers(ctx context.Context, filtersQuery request_models.PgnFltQueryParams) ([]dto.ResourceServer, *model.MetaData, error)
}

type MiniRideModule interface {
	ListenMiniRideEvent(ctx context.Context)
	ProcessEvents(ctx context.Context, miniRideEvent *request_models.MinRideEvent, wg *sync.WaitGroup)
	CheckPhone(ctx context.Context, phone string) (*dto.MiniRideResponse, error)
}
