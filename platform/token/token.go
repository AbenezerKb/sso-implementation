package token

import (
	"context"
	"crypto/rsa"
	"fmt"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/platform"
	"sso/platform/logger"
	"sso/platform/utils"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

type Jwt struct {
	logger     logger.Logger
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func JwtInit(logger logger.Logger, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) platform.Token {
	return &Jwt{
		logger:     logger,
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (j *Jwt) GenerateAccessToken(ctx context.Context, userID string, expiresAt time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
		Issuer:    "test",
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   userID,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodPS512, claims).SignedString(j.privateKey)
	if err != nil {
		j.logger.Error(ctx, "could not generate access token", zap.Error(err))
		return "", errors.ErrInternalServerError.Wrap(err, "could not generate access token")
	}
	return token, nil
}
func (j *Jwt) GenerateAccessTokenForClient(ctx context.Context, userID, clientID, scope string, expiresAt time.Duration) (string, error) {
	claims := dto.AccessToken{
		Scope: scope,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
			Issuer:    "test",
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   userID,
			Audience:  jwt.ClaimStrings{clientID},
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodPS512, claims).SignedString(j.privateKey)
	if err != nil {
		j.logger.Error(ctx, "could not generate access token", zap.Error(err))
		return "", errors.ErrInternalServerError.Wrap(err, "could not generate access token")
	}
	return token, nil
}

func (j *Jwt) GenerateRefreshToken(ctx context.Context) string {
	return utils.GenerateRandomString(25, true)
}

func (j *Jwt) GenerateIdToken(ctx context.Context, user *dto.User, clientId string, expiresAt time.Duration) (string, error) {
	claims := dto.IDTokenPayload{
		FirstName:   user.FirstName,
		MiddleName:  user.MiddleName,
		LastName:    user.LastName,
		Picture:     user.ProfilePicture,
		Email:       user.Email,
		PhoneNumber: user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			Audience:  jwt.ClaimStrings{clientId},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodPS512, claims).SignedString(j.privateKey)
	if err != nil {
		j.logger.Error(ctx, "could not generate id token", zap.Error(err))
		return "", errors.ErrInternalServerError.Wrap(err, "could not generate id token")
	}
	return token, nil
}

func (j *Jwt) VerifyToken(signingMethod jwt.SigningMethod, token string) (bool, *jwt.RegisteredClaims) {
	claims := &jwt.RegisteredClaims{}

	segments := strings.Split(token, ".")
	if len(segments) < 3 {
		return false, claims
	}
	err := signingMethod.Verify(strings.Join(segments[:2], "."), segments[2], j.publicKey)
	if err != nil {
		return false, claims
	}

	if _, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSAPSS); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return j.publicKey, nil
	}); err != nil {
		return false, claims
	}
	return true, claims
}

func (j *Jwt) VerifyIdToken(signingMethod jwt.SigningMethod, token string) (bool, *dto.IDTokenPayload) {
	claims := &dto.IDTokenPayload{}

	segments := strings.Split(token, ".")
	if len(segments) < 3 {
		return false, claims
	}

	err := signingMethod.Verify(strings.Join(segments[:2], "."), segments[2], j.publicKey)
	if err != nil {
		return false, claims
	}

	if _, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSAPSS); !ok {
			return nil, fmt.Errorf("unexpected siging method %v", t.Header["alg"])
		}
		return j.publicKey, nil
	}); err != nil {
		return false, claims
	}

	return true, claims

}
