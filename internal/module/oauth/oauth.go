package oauth

import (
	"context"
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
}

func SetOptions(options Options) Options {
	if options.AccessTokenExpireTime == 0 {
		options.AccessTokenExpireTime = time.Minute * 10
	}
	if options.RefreshTokenExpireTime == 0 {
		options.RefreshTokenExpireTime = time.Hour * 24 * 30
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
	refreshToken := o.token.GenerateRefreshToken(ctx)

	// TODO: persist the refresh token
	err = o.oauthPersistence.SaveInternalRefreshToken(ctx, dto.InternalRefreshToken{
		Refreshtoken: refreshToken,
		UserID:       user.ID,
	})
	if err != nil {
		return nil, err
	}

	idToken, err := o.token.GenerateIdToken(ctx, user)
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

func (o *oauth) Logout(ctx context.Context) error {
	id, ok := ctx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid user id")
		o.logger.Error(ctx, "invalid user id", zap.Error(err), zap.Any("user_id", id))
		return err
	}
	userId, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "could not parse user id")
		o.logger.Error(ctx, "parse error", zap.Error(err))
		return err
	}

	err = o.oauthPersistence.RemoveInternalRefreshToken(ctx, userId)
	if err != nil {
		return err
	}

	return nil
}
