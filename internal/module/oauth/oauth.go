package oauth

import (
	"context"
	"fmt"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform"
	"sso/platform/logger"
	"sso/platform/utils"
	"time"

	"github.com/joomcode/errorx"

	"github.com/dongri/phonenumber"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type oauth struct {
	logger           logger.Logger
	oauthPersistence storage.OAuthPersistence
	ipPersistence    storage.IdentityProviderPersistence
	otpCache         storage.OTPCache
	sessionCache     storage.SessionCache
	token            platform.Token
	smsClient        platform.SMSClient
	options          Options
	selfIP           platform.IdentityProvider
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
func InitOAuth(logger logger.Logger,
	oauthPersistence storage.OAuthPersistence,
	ipPersistence storage.IdentityProviderPersistence,
	otpCache storage.OTPCache,
	sessionCache storage.SessionCache,
	token platform.Token,
	smsClient platform.SMSClient,
	selfIP platform.IdentityProvider,
	options Options) module.OAuthModule {
	return &oauth{
		logger:           logger,
		oauthPersistence: oauthPersistence,
		ipPersistence:    ipPersistence,
		otpCache:         otpCache,
		sessionCache:     sessionCache,
		token:            token,
		smsClient:        smsClient,
		selfIP:           selfIP,
		options:          options,
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

func (o *oauth) Login(ctx context.Context, userParam dto.LoginCredential, userDeviceAddress dto.UserDeviceAddress) (*dto.TokenResponse, error) {
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

	if user.Status != constant.Active {
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
	err = o.oauthPersistence.SaveInternalRefreshToken(ctx, dto.InternalRefreshToken{
		RefreshToken: refreshToken,
		UserID:       user.ID,
		UserAgent:    userDeviceAddress.UserAgent,
		IPAddress:    userDeviceAddress.IPAddress,
		ExpiresAt:    time.Now().Add(o.options.RefreshTokenExpireTime),
	})
	if err != nil {
		return nil, err
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

	if err := o.oauthPersistence.RemoveInternalRefreshToken(ctx, oldRefreshToken.RefreshToken); err != nil {
		return err
	}

	return nil
}

func (o *oauth) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	oldRefreshToken, err := o.oauthPersistence.GetInternalRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if time.Now().After(oldRefreshToken.ExpiresAt) {
		if err := o.oauthPersistence.RemoveInternalRefreshToken(ctx, oldRefreshToken.RefreshToken); err != nil {
			return nil, err
		}

		err := errors.ErrAuthError.New("internal refresh token expired")
		o.logger.Warn(ctx, "internal token expired", zap.Error(err), zap.String("internal refresh token", oldRefreshToken.RefreshToken))
		return nil, err
	}

	accessToken, err := o.token.GenerateAccessToken(ctx, oldRefreshToken.UserID.String(), o.options.AccessTokenExpireTime)
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

	newRefreshToken, err := o.oauthPersistence.UpdateInternalRefreshToken(ctx, oldRefreshToken.RefreshToken, o.token.GenerateRefreshToken(ctx))
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.RefreshToken,
		TokenType:    constant.BearerToken,
		IDToken:      idToken,
		ExpiresIn:    fmt.Sprintf("%vs", o.options.AccessTokenExpireTime.Seconds()),
	}, nil
}

func (o *oauth) LoginWithIdentityProvider(ctx context.Context, login request_models.LoginWithIP, userDeviceAddress dto.UserDeviceAddress) (dto.TokenResponse, error) {
	// validate
	if err := login.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		o.logger.Info(ctx, "invalid input on login with identity provider", zap.Error(err), zap.Any("login", login))
		return dto.TokenResponse{}, err
	}
	// check and get if ip exists
	ipID, err := uuid.Parse(login.IdentityProvider)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid identity provider")
		o.logger.Info(ctx, "invalid identity provider id", zap.Error(err), zap.String("ip-id", login.IdentityProvider))
		return dto.TokenResponse{}, err
	}
	ip, err := o.ipPersistence.GetIdentityProvider(ctx, ipID)
	if err != nil {
		if errorx.IsOfType(err, errors.ErrNoRecordFound) {
			err := errors.ErrInvalidUserInput.Wrap(err, fmt.Sprintf("identity provider with id %s does not exist", ipID.String()))
			return dto.TokenResponse{}, err
		}
		return dto.TokenResponse{}, err
	}
	// request platform
	// FixMe: decrypt client secret
	accessToken, refreshToken, err := o.selfIP.GetAccessToken(ctx, ip.TokenEndpointURI, ip.RedirectURI, ip.ClientID, ip.ClientSecret, login.Code)
	if err != nil {
		err := errors.ErrAuthError.Wrap(err, "authentication failed")
		o.logger.Info(ctx, "login authentication for identity provider failed", zap.Error(err), zap.Any("login", login), zap.Any("ip", ip))
		return dto.TokenResponse{}, err
	}
	userInfo, err := o.selfIP.GetUserInfo(ctx, ip.UserInfoEndpointURI, accessToken)
	if err != nil {
		err := errors.ErrAcessError.Wrap(err, "authorization for user-info failed")
		o.logger.Warn(ctx, "getting user info from identity provider failed", zap.Error(err), zap.Any("ip", ip))
		return dto.TokenResponse{}, err
	}
	if err := userInfo.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid userinfo")
		o.logger.Warn(ctx, "invalid userinfo was returned from identity provider", zap.Any("user-info", userInfo), zap.Error(err))
		return dto.TokenResponse{}, err
	}
	// save or update access token
	ipAT, err := o.ipPersistence.GetIPAccessTokenBySubAndIP(ctx, userInfo.Sub, ip.ID)
	var user *dto.User
	if err != nil {
		if !errorx.IsOfType(err, errors.ErrNoRecordFound) {
			return dto.TokenResponse{}, err
		}
		// check for email uniqueness
		exists, err := o.oauthPersistence.UserByEmailExists(ctx, userInfo.Email)
		if err != nil {
			return dto.TokenResponse{}, err
		}
		if exists {
			err := errors.ErrInvalidUserInput.Wrap(err, "user with this email already exists")
			o.logger.Info(ctx, "user with email already exists for login with ip", zap.Error(err), zap.String("email", userInfo.Email))
			return dto.TokenResponse{}, err
		}
		// check for phone uniqueness
		exists, err = o.oauthPersistence.UserByPhoneExists(ctx, userInfo.Phone)
		if err != nil {
			return dto.TokenResponse{}, err
		}
		if exists {
			err := errors.ErrInvalidUserInput.Wrap(err, "user with this phone already exists")
			o.logger.Info(ctx, "user with phone already exists for login with ip", zap.Error(err), zap.String("phone", userInfo.Phone))
			return dto.TokenResponse{}, err
		}
		// save user
		user, err = o.oauthPersistence.Register(ctx, dto.User{
			FirstName:      userInfo.FirstName,
			MiddleName:     userInfo.MiddleName,
			LastName:       userInfo.LastName,
			Email:          userInfo.Email,
			Phone:          userInfo.Phone,
			Gender:         userInfo.Gender,
			ProfilePicture: userInfo.ProfilePicture,
		})
		if err != nil {
			return dto.TokenResponse{}, err
		}
		// create access token
		ipAT, err = o.ipPersistence.SaveIPAccessToken(ctx, dto.IPAccessToken{
			UserID:       user.ID,
			SubID:        userInfo.Sub,
			IPID:         ip.ID,
			Token:        accessToken,
			RefreshToken: refreshToken,
		})
		if err != nil {
			return dto.TokenResponse{}, err
		}
	} else {
		// update access token
		ipAT.Token = accessToken
		ipAT.RefreshToken = refreshToken
		ipAT, err = o.ipPersistence.UpdateIpAccessToken(ctx, ipAT)
		if err != nil {
			return dto.TokenResponse{}, err
		}
		// get user
		user, err = o.oauthPersistence.GetUserByID(ctx, ipAT.UserID)
		if err != nil {
			return dto.TokenResponse{}, err
		}
	}

	internalAccessToken, err := o.token.GenerateAccessToken(ctx, user.ID.String(), o.options.AccessTokenExpireTime)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	internalRefreshToken := o.token.GenerateRefreshToken(ctx)

	err = o.oauthPersistence.SaveInternalRefreshToken(ctx, dto.InternalRefreshToken{
		RefreshToken: internalRefreshToken,
		UserID:       user.ID,
		UserAgent:    userDeviceAddress.UserAgent,
		IPAddress:    userDeviceAddress.IPAddress,
		ExpiresAt:    time.Now().Add(o.options.RefreshTokenExpireTime),
	})

	if err != nil {
		return dto.TokenResponse{}, err
	}

	idToken, err := o.token.GenerateIdToken(ctx, user, "sso", o.options.IDTokenExpireTime)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	accessTokenResponse := dto.TokenResponse{
		AccessToken:  internalAccessToken,
		RefreshToken: internalRefreshToken,
		IDToken:      idToken,
		TokenType:    constant.BearerToken,
		ExpiresIn:    fmt.Sprintf("%vs", o.options.AccessTokenExpireTime.Seconds()),
	}
	return accessTokenResponse, nil
}

func (o *oauth) GetAllIdentityProviders(ctx context.Context) ([]dto.IdentityProvider, error) {
	return o.oauthPersistence.GetAllIdentityProviders(ctx)
}
