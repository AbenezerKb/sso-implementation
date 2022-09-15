package oauth

import (
	"context"
	"fmt"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform"
	"sso/platform/logger"
	"sso/platform/utils"
	"time"

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
	token            platform.Token
	smsClient        platform.SMSClient
	options          Options
}

type Options struct {
	AccessTokenExpireTime  time.Duration
	RefreshTokenExpireTime time.Duration
	IDTokenExpireTime      time.Duration
}

func SetOptions(options Options) Options {
	if options.AccessTokenExpireTime == 0 {
		options.AccessTokenExpireTime = time.Minute * 10
	}
	if options.RefreshTokenExpireTime == 0 {
		options.RefreshTokenExpireTime = time.Hour * 24 * 30
	}
	if options.IDTokenExpireTime == 0 {
		options.IDTokenExpireTime = time.Minute * 10
	}
	return options
}
func InitOAuth(logger logger.Logger, oauthPersistence storage.OAuthPersistence, otpCache storage.OTPCache, sessionCache storage.SessionCache, token platform.Token, smsClient platform.SMSClient, options Options) module.OAuthModule {
	return &oauth{
		logger,
		oauthPersistence,
		otpCache,
		sessionCache,
		token,
		smsClient,
		options,
	}
}

func (o *oauth) Register(ctx context.Context, userParam dto.RegisterUser) (*dto.User, error) {
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

	userParam.Password, err = utils.HashAndSalt(ctx, []byte(userParam.Password), o.logger)
	if err != nil {
		return nil, err
	}

	user, err := o.oauthPersistence.Register(ctx, userParam.User)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (o *oauth) Login(ctx context.Context, userParam dto.LoginCredential) (*dto.TokenResponse, error) {
	if err := userParam.ValidateLoginCredential(); err != nil {
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

	accessToken, err := o.token.GenerateAccessToken(ctx, user.ID.String(), o.options.AccessTokenExpireTime)
	if err != nil {
		return nil, err
	}
	oldRfToken, err := o.oauthPersistence.GetInternalRefreshTokenByUserID(ctx, user.ID)
	refreshToken := ""
	if err != nil {
		refreshToken = o.token.GenerateRefreshToken(ctx)

		err = o.oauthPersistence.SaveInternalRefreshToken(ctx, dto.InternalRefreshToken{
			Refreshtoken: refreshToken,
			UserID:       user.ID,
			ExpiresAt:    time.Now().Add(o.options.RefreshTokenExpireTime),
		})

		if err != nil {
			return nil, err
		}
	} else if time.Now().After(oldRfToken.ExpiresAt) {
		if err := o.oauthPersistence.RemoveInternalRefreshToken(ctx, oldRfToken.Refreshtoken); err != nil {
			return nil, err
		}
		refreshToken = o.token.GenerateRefreshToken(ctx)

		err = o.oauthPersistence.SaveInternalRefreshToken(ctx, dto.InternalRefreshToken{
			Refreshtoken: refreshToken,
			UserID:       user.ID,
			ExpiresAt:    time.Now().Add(o.options.RefreshTokenExpireTime),
		})

		if err != nil {
			return nil, err
		}

	} else {
		refreshToken = oldRfToken.Refreshtoken
	}

	idToken, err := o.token.GenerateIdToken(ctx, user, "sso", o.options.IDTokenExpireTime)
	if err != nil {
		return nil, err
	}

	accessTokenResponse := dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IDToken:      idToken,
		TokenType:    constant.BearerToken,
		ExpiresIn:    fmt.Sprintf("%vs", o.options.AccessTokenExpireTime.Seconds()),
	}
	return &accessTokenResponse, nil
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

	if user.Status != constant.Active {
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

func (o *oauth) Logout(ctx context.Context, param dto.InternalRefreshTokenRequestBody) error {
	if err := param.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil
	}
	oldRefreshToken, err := o.oauthPersistence.GetInternalRefreshToken(ctx, param.RefreshToken)
	if err != nil {
		return err
	}

	if err := o.oauthPersistence.RemoveInternalRefreshToken(ctx, oldRefreshToken.Refreshtoken); err != nil {
		return err
	}

	return nil
}

func (o *oauth) RefreshToken(ctx context.Context, param dto.InternalRefreshTokenRequestBody) (*dto.TokenResponse, error) {
	if err := param.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	oldRefreshToken, err := o.oauthPersistence.GetInternalRefreshToken(ctx, param.RefreshToken)
	if err != nil {
		return nil, err
	}

	if time.Now().After(oldRefreshToken.ExpiresAt) {
		if err := o.oauthPersistence.RemoveInternalRefreshToken(ctx, oldRefreshToken.Refreshtoken); err != nil {
			return nil, err
		}

		err := errors.ErrAuthError.New("internal refresh token expired")
		o.logger.Warn(ctx, "internal token expired", zap.Error(err), zap.String("internal refresh token", oldRefreshToken.Refreshtoken))
		return nil, err
	}

	accessToken, err := o.token.GenerateAccessToken(ctx, oldRefreshToken.UserID.String(), o.options.AccessTokenExpireTime)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	user, err := o.oauthPersistence.GetUserByID(ctx, oldRefreshToken.UserID)
	if err != nil {
		return nil, err
	}
	idToken, err := o.token.GenerateIdToken(ctx, user, "sso", o.options.IDTokenExpireTime)
	if err != nil {
		return nil, err
	}
	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: oldRefreshToken.Refreshtoken,
		TokenType:    constant.BearerToken,
		IDToken:      idToken,
		ExpiresIn:    fmt.Sprintf("%vs", o.options.AccessTokenExpireTime.Seconds()),
	}, nil
}
