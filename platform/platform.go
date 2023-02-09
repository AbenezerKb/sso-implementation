package platform

import (
	"context"
	"mime/multipart"
	"time"

	"sso/internal/constant/model/dto"

	"github.com/golang-jwt/jwt/v4"
)

type SMSConfig struct {
	UserName  string
	Password  string
	Sender    string
	DLRMask   string
	DCS       string
	DLRURL    string
	Server    string
	Templates map[string]string
	Type      string
	APIKey    string
}

type SMSClient interface {
	SendSMS(ctx context.Context, to, text string) error
	SendSMSWithTemplate(ctx context.Context, to, templateName string, values ...interface{}) error
}

type Token interface {
	GenerateAccessToken(ctx context.Context, userID string, expiresAt time.Duration) (string, error)
	GenerateAccessTokenForClient(ctx context.Context, userID, clientID, scope string, expiresAt time.Duration) (string, error)
	GenerateRefreshToken(ctx context.Context) string
	GenerateIdToken(ctx context.Context, user *dto.User, clientId string, expiresAt time.Duration) (string, error)
	VerifyToken(signingMethod jwt.SigningMethod, token string) (bool, *jwt.RegisteredClaims)
	VerifyIdToken(signingMethod jwt.SigningMethod, token string) (bool, *dto.IDTokenPayload)
}

type IdentityProvider interface {
	GetAccessToken(ctx context.Context, endPoint, redirectURI, clientID, clientSecret, code string) (string, string, error)
	GetUserInfo(ctx context.Context, endPoint, accessToken string) (dto.UserInfo, error)
}

type Asset interface {
	SaveAsset(ctx context.Context, asset multipart.File, dst string) error
}
