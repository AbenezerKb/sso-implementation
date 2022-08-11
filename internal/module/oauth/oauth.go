package oauth

import (
	"context"
	"crypto/rsa"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform"
	"sso/platform/logger"

	"github.com/dongri/phonenumber"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type oauth struct {
	logger           logger.Logger
	oauthPersistence storage.OAuthPersistence
	otpCache         storage.OTPCache
	sessionCache     storage.SessionCache
	privateKey       *rsa.PrivateKey
	smsClient        platform.SMSClient
}

func InitOAuth(logger logger.Logger, oauthPersistence storage.OAuthPersistence, otpCache storage.OTPCache, sessionCache storage.SessionCache, privateKey *rsa.PrivateKey, smsClient platform.SMSClient) module.OAuthModule {
	return &oauth{
		logger,
		oauthPersistence,
		otpCache,
		sessionCache,
		privateKey,
		smsClient,
	}
}

func (o *oauth) Register(ctx context.Context, userParam dto.User) (*dto.User, error) {
	if err := userParam.ValidateUser(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}
	userParam.Phone = phonenumber.Parse(userParam.Phone, "ET")

	err := o.VerifyOTP(ctx, userParam.Phone, userParam.OTP)
	if err != nil {
		return nil, err
	}

	exists, err := o.oauthPersistence.UserByPhoneExists(ctx, userParam.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		o.logger.Info(ctx, "user already exists", zap.String("phone", userParam.Phone))
		return nil, errors.ErrDataExists.New("user with this phone already exists")
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

	if userParam.Email != "" && userParam.Password != "" {
		query = userParam.Email
	} else if userParam.Phone != "" && userParam.OTP != "" {
		userParam.Phone = phonenumber.Parse(userParam.Phone, "ET")
		query = userParam.Phone
	}

	user, err := o.oauthPersistence.GetUserByPhoneOrEmail(ctx, query)

	if err != nil {
		return nil, errors.ErrInvalidUserInput.Wrap(err, "invalid credentials")
	}

	if user.Status != "ACTIVE" {
		err := errors.ErrInvalidUserInput.New("Account is deactivated")
		o.logger.Info(ctx, "user is not active", zap.Error(err))
		return nil, err
	}

	if userParam.Email != "" && userParam.Password != "" {
		if !o.ComparePassword(user.Password, userParam.Password) {
			err := errors.ErrInvalidUserInput.New("Invalid credentials")
			o.logger.Info(ctx, "invalid credentials", zap.Error(err))
			return nil, err
		}
	} else if userParam.Phone != "" && userParam.OTP != "" {
		err := o.VerifyOTP(ctx, userParam.Phone, userParam.OTP)
		if err != nil {
			return nil, err
		}

	}

	accessToken, err := o.GenerateAccessToken(ctx, user)
	if err != nil {
		return nil, err
	}
	refreshToken, err := o.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}
	// TODO: persist the refresh token
	//err = o.cache.Set(ctx, refreshToken, user.ID.String(), time.Hour*24*7).Err()
	//if err != nil {
	//	o.logger.Error(ctx, "could not persist refresh token", zap.Error(err))
	//	return nil, errors.ErrCacheSetError.Wrap(err, "could not persist refresh token")
	//}

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

func (o *oauth) VerifyUserStatus(ctx context.Context, phone string) error {
	user, err := o.oauthPersistence.GetUserByPhone(ctx, phone)
	if err != nil {
		return err
	}

	if user.Status != "ACTIVE" {
		err := errors.ErrInvalidUserInput.New("Account is deactivated")
		o.logger.Info(ctx, "user is not active", zap.Error(err))
		return err
	}
	return nil
}
func (o *oauth) GetUserStatus(ctx context.Context, Id string) (string, error) {
	userId, err := uuid.Parse(Id)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "could not parse user id")
		o.logger.Error(ctx, "parse error", zap.Error(err))
		return "", err
	}
	status, err := o.oauthPersistence.GetUserStatus(ctx, userId)
	if err != nil {
		return "", err
	}

	return status, nil
}
