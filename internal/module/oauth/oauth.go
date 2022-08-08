package oauth

import (
	"context"
	"crypto/rsa"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dongri/phonenumber"
	"github.com/go-redis/redis/v8"
	"github.com/joomcode/errorx"
	"go.uber.org/zap"
)

type oauth struct {
	logger           logger.Logger
	oauthPersistence storage.OAuthPersistence
	cache            *redis.Client
	privateKey       *rsa.PrivateKey
}

func InitOAuth(logger logger.Logger, oauthPersistence storage.OAuthPersistence, cache *redis.Client, privateKey *rsa.PrivateKey) module.OAuthModule {
	return &oauth{
		logger,
		oauthPersistence,
		cache,
		privateKey,
	}
}

func (o *oauth) Register(ctx context.Context, userParam dto.User) (*dto.User, error) {
	if err := userParam.ValidateUser(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	exists, err := o.oauthPersistence.UserByPhoneExists(ctx, userParam.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrDataExists.Wrap(err, "user with this phone already exists")
	}

	if userParam.Email != "" {
		exists, err := o.oauthPersistence.UserByEmailExists(ctx, userParam.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.ErrDataExists.Wrap(err, "user with this email already exists")
		}
	}

	userParam.Password, err = o.HashAndSalt(ctx, []byte(userParam.Password))
	if err != nil {
		return nil, err
	}

	user, err := o.oauthPersistence.Register(ctx, userParam)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (o *oauth) Login(ctx context.Context, userParam dto.User) (*dto.TokenResponse, error) {
	if err := userParam.ValidateLoginCredentials(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	var query string

	if userParam.UserName != "" && userParam.Password != "" {
		query = userParam.UserName
	} else if userParam.Phone != "" && userParam.OTP != "" {
		userParam.Phone = phonenumber.Parse(userParam.Phone, "ET")
		query = userParam.Phone
	} else {
		err := errors.ErrInvalidUserInput.Wrap(&errorx.Error{}, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	// user, err := o.oauthPersistence.GetUserByPhoneOrUserNameOrEmail(ctx, query)
	user, err := o.oauthPersistence.GetUserByPhone(ctx, query)

	if err != nil {
		return nil, errors.ErrNoRecordFound.Wrap(err, "user not found")
	}

	if user.Status != "ACTIVE" {
		// Todo: add error Account is deactivated
		return nil, errors.ErrNoRecordFound.Wrap(err, "Account is deactivated")
	}

	if user.UserName != "" && user.Password != "" {
		if !o.ComparePassword(user.Password, userParam.Password) {
			return nil, errors.ErrInvalidUserInput.Wrap(err, "invalid password")
		}
	} else if user.Phone != "" && user.OTP != "" {
		// Todo: verify OTP
	}

	accessToken, err := o.GenerateAccessToken(ctx, user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := o.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}
	// persist the refresh token
	err = o.cache.Set(ctx, refreshToken, user.ID.String(), time.Hour*24*7).Err()
	if err != nil {
		return nil, err
	}

	idToken, err := o.GenerateIdToken(ctx, user)
	if err != nil {
		return nil, err
	}

	accessTokenResponse := dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IDToken:      idToken,
	}
	return &accessTokenResponse, nil
}
func (o *oauth) HashAndSalt(ctx context.Context, pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		o.logger.Error(ctx, "could not hash password", zap.Error(err))
		return "", err
	}
	return string(hash), nil
}

func (o *oauth) ComparePassword(hashedPwd, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPassword))
	return err == nil
}
