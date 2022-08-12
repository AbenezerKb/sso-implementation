package oauth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/golang-jwt/jwt/v4"
)

func (o *oauth) GenerateAccessToken(ctx context.Context, user *dto.User) (string, error) {
	claims := dto.AccessToken{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(o.options.AccessTokenExpireTime)),
			Issuer:    "test",
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodPS512, claims).SignedString(o.privateKey)
	if err != nil {
		o.logger.Error(ctx, "could not generate access token", zap.Error(err))
		return "", errors.ErrInternalServerError.Wrap(err, "could not generate access token")
	}
	return token, nil
}

func (o *oauth) GenerateRefreshToken(ctx context.Context, user *dto.User) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(o.options.RefreshTokenExpireTime)),
		Issuer:    "test",
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   user.ID.String(),
	}

	rfToken, err := jwt.NewWithClaims(jwt.SigningMethodPS512, claims).SignedString(o.privateKey)
	if err != nil {
		o.logger.Error(ctx, "could not generate refresh token", zap.Error(err))
		return "", errors.ErrInternalServerError.Wrap(err, "could not generate refresh token")
	}

	return rfToken, nil
}

func (o *oauth) GenerateIdToken(ctx context.Context, user *dto.User) (string, error) {
	claims := dto.IDTokenPayload{
		FirstName:   user.FirstName,
		MiddleName:  user.MiddleName,
		LastName:    user.LastName,
		Picture:     user.ProfilePicture,
		Email:       user.Email,
		PhoneNumber: user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(user.CreatedAt),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodPS512, claims).SignedString(o.privateKey)
	if err != nil {
		o.logger.Error(ctx, "could not generate id token", zap.Error(err))
		return "", errors.ErrInternalServerError.Wrap(err, "could not generate id token")
	}
	return token, nil
}

func (o *oauth) VerifyToken(signingMethod jwt.SigningMethod, token string, pk *rsa.PublicKey) (bool, *jwt.RegisteredClaims) {
	claims := &jwt.RegisteredClaims{}

	segments := strings.Split(token, ".")
	if len(segments) < 3 {
		return false, claims
	}
	err := signingMethod.Verify(strings.Join(segments[:2], "."), segments[2], pk)
	if err != nil {
		return false, claims
	}

	if _, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSAPSS); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return pk, nil
	}); err != nil {
		return false, claims
	}
	return true, claims
}
