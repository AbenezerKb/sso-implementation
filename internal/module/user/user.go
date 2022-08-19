package user

import (
	"context"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform"
	"sso/platform/logger"
	"sso/platform/utils"

	"github.com/casbin/casbin/v2"
	"github.com/dongri/phonenumber"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	logger           logger.Logger
	oauthPersistence storage.OAuthPersistence
	smsClient        platform.SMSClient
	enforcer         *casbin.Enforcer
}

func Init(logger logger.Logger, oauthPersistence storage.OAuthPersistence, smsClient platform.SMSClient, enforcer *casbin.Enforcer) module.UserModule {
	return &user{
		logger,
		oauthPersistence,
		smsClient,
		enforcer,
	}
}

func (u *user) Create(ctx context.Context, param dto.CreateUser) (*dto.User, error) {
	if err := param.ValidateUser(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		u.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	param.Phone = phonenumber.Parse(param.Phone, "ET")
	exists, err := u.oauthPersistence.UserByPhoneExists(ctx, param.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		u.logger.Info(ctx, "user already exists", zap.String("phone", param.Phone))
		return nil, errors.ErrDataExists.New("user with this phone already exists")
	}

	exists, err = u.oauthPersistence.UserByEmailExists(ctx, param.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		u.logger.Info(ctx, "user already exists", zap.String("email", param.Email))
		return nil, errors.ErrDataExists.Wrap(err, "user with this email already exists")
	}

	password := utils.GenerateRandomString(6, false)

	param.Password, err = utils.HashAndSalt(ctx, []byte(password), u.logger)
	if err != nil {
		return nil, err
	}
	err = u.smsClient.SendSMSWithTemplate(ctx, param.Phone, "password", string(password))
	if err != nil {
		return nil, err
	}
	user, err := u.oauthPersistence.Register(ctx, param.User)
	if err != nil {
		return nil, err
	}
	if exists, _ := u.enforcer.HasRoleForUser(param.Role, user.ID.String()); !exists {
		u.enforcer.AddRoleForUser(user.ID.String(), param.Role)
	}
	return user, nil
}
func (u *user) HashAndSalt(ctx context.Context, pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, 14)
	if err != nil {
		u.logger.Error(ctx, "could not hash password", zap.Error(err))
		return "", err
	}
	return string(hash), nil
}
