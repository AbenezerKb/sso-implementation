package platform

import (
	"context"
	"sso/internal/constant/model/dto"
	"time"

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
	GenerateRefreshToken(ctx context.Context) string
	GenerateIdToken(ctx context.Context, user *dto.User, clientId string, expiresAt time.Duration) (string, error)
	VerifyToken(signingMethod jwt.SigningMethod, token string) (bool, *jwt.RegisteredClaims)
	VerifyIdToken(signingMethod jwt.SigningMethod, token string) (bool, *dto.IDTokenPayload)
}
