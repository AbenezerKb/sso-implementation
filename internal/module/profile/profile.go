package profile

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/module"
	"sso/internal/storage"
	"sso/platform/logger"
	"sso/platform/utils"
	"strings"
	"time"

	"github.com/dongri/phonenumber"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Options struct {
	PublicAddress         string
	ProfilePictureDist    string
	ProfilePictureMaxSize int
}

func SetOptions(options Options) Options {
	if options.PublicAddress == "" {
		options.PublicAddress = "http://127.0.0.1:8000/assets/"
	}

	if options.ProfilePictureDist == "" {
		options.ProfilePictureDist = "../../../../static/profile_picture/"
	}

	if options.ProfilePictureMaxSize == 0 {
		options.ProfilePictureMaxSize = 2000001
	}

	return options
}

type profileModule struct {
	logger             logger.Logger
	oauthPersistence   storage.OAuthPersistence
	profilePersistence storage.ProfilePersistence
	otpCache           storage.OTPCache
	options            Options
}

func InitProfile(logger logger.Logger, oauthPersistence storage.OAuthPersistence, profilePersistence storage.ProfilePersistence, otpCache storage.OTPCache, options Options) module.ProfileModule {
	return &profileModule{
		logger:             logger,
		oauthPersistence:   oauthPersistence,
		profilePersistence: profilePersistence,
		otpCache:           otpCache,
		options:            options,
	}
}

func (p *profileModule) UpdateProfile(ctx context.Context, userParam dto.User) (*dto.User, error) {
	id, ok := ctx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid user id")
		p.logger.Info(ctx, "invalid user id", zap.Error(err), zap.Any("user_id", id))
		return &dto.User{}, err
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "user not found")
		p.logger.Info(ctx, "parse error", zap.Error(err), zap.String("user id", id))
		return nil, err
	}

	if err := userParam.ValidateUpdateProfile(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.logger.Info(ctx, "invalid input", zap.Error(err))
		return nil, err
	}

	userParam.ID = userID
	updatedUser, err := p.profilePersistence.UpdateProfile(ctx, userParam)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (p *profileModule) GetProfile(ctx context.Context) (*dto.User, error) {
	id, ok := ctx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid user id")
		p.logger.Info(ctx, "invalid user id", zap.Error(err), zap.Any("user_id", id))
		return &dto.User{}, err
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "user not found")
		p.logger.Info(ctx, "parse error", zap.Error(err), zap.String("user id", id))
		return nil, err
	}
	user, err := p.profilePersistence.GetProfile(ctx, userID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *profileModule) UpdateProfilePicture(ctx context.Context, imageFile *multipart.FileHeader) error {
	id, ok := ctx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid user id")
		p.logger.Info(ctx, "invalid user id", zap.Error(err), zap.Any("user_id", id))
		return err
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "user not found")
		p.logger.Info(ctx, "parse error", zap.Error(err), zap.String("user id", id))
		return err
	}

	src, err := imageFile.Open()
	if err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid picture")
		p.logger.Info(ctx, "couldn't read image", zap.Error(err), zap.Any("image", imageFile), zap.Any("user-id", id))
		return err
	}
	defer src.Close()

	buff := make([]byte, 512)
	if _, err := src.Read(buff); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid picture")
		p.logger.Info(ctx, "couldn't read image", zap.Error(err), zap.Any("image", imageFile), zap.Any("user-id", id))
		return err
	}
	fileType := http.DetectContentType(buff)

	if strings.Split(fileType, "/")[0] != "image" {
		err = errors.ErrInvalidUserInput.New("invalid picture")
		p.logger.Info(ctx, "provided profile picture is not image", zap.Error(err), zap.Any("image", imageFile), zap.Any("user-id", id))
		return err
	}

	if imageFile.Size > int64(p.options.ProfilePictureMaxSize) {
		err = errors.ErrInvalidUserInput.New("image size must be less than 2MB")
		p.logger.Info(ctx, "image size too big", zap.Error(err), zap.String("image", imageFile.Filename), zap.Any("size", imageFile.Size), zap.Any("user-id", id))
		return err
	}

	// final image name
	finalImageName := fmt.Sprint(time.Now().UnixMilli()) + "_" + id + "_" + imageFile.Filename

	err = utils.SaveMultiPartFile(imageFile, p.options.ProfilePictureDist+"/"+finalImageName)
	if err != nil {
		err = errors.ErrInternalServerError.Wrap(err, "couldn't save profile picture")
		p.logger.Error(context.Background(), "error unable to save profile picture to disck", zap.Error(err), zap.Any("image", imageFile))
	}

	err = p.profilePersistence.UpdateProfilePicture(ctx, p.options.PublicAddress+"/"+finalImageName, userID)
	if err != nil {
		return err
	}

	return nil
}

func (p *profileModule) ChangePhone(ctx context.Context, changePhoneParam dto.ChangePhoneParam) error {
	id, ok := ctx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid user id")
		p.logger.Info(ctx, "invalid user id", zap.Error(err), zap.Any("user_id", id))
		return err
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "user not found")
		p.logger.Info(ctx, "parse error", zap.Error(err), zap.String("user id", id))
		return err
	}

	if err := changePhoneParam.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.logger.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	changePhoneParam.Phone = phonenumber.Parse(changePhoneParam.Phone, "ET")

	err = p.otpCache.VerifyOTP(ctx, changePhoneParam.Phone, changePhoneParam.OTP)
	if err != nil {
		return err
	}

	exists, err := p.oauthPersistence.UserByPhoneExists(ctx, changePhoneParam.Phone)
	if err != nil {
		return err
	}
	if exists {
		p.logger.Info(ctx, "user already exists", zap.String("phone", changePhoneParam.Phone))
		return errors.ErrDataExists.New("user with this phone already exists")
	}

	return p.profilePersistence.ChangePhone(ctx, changePhoneParam, userID)

}

func (p *profileModule) ChangePassword(ctx context.Context, changePasswordParam dto.ChangePasswordParam) error {
	id, ok := ctx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid user id")
		p.logger.Info(ctx, "invalid user id", zap.Error(err), zap.Any("user_id", id))
		return err
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "user not found")
		p.logger.Info(ctx, "parse error", zap.Error(err), zap.String("user id", id))
		return err
	}

	if err := changePasswordParam.Validate(); err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.logger.Info(ctx, "invalid input", zap.Error(err))
		return err
	}

	userPassword, err := p.oauthPersistence.GetUserPassword(ctx, userID)
	if err != nil {
		return err
	}

	if !utils.CompareHashAndPassword(userPassword, changePasswordParam.OldPassword) {
		err := errors.ErrInvalidUserInput.New("invalid credentials")
		p.logger.Info(ctx, "invalid credentials", zap.Error(err))
		return err
	}

	changePasswordParam.NewPassword, err = utils.HashAndSalt(ctx, []byte(changePasswordParam.NewPassword), p.logger)
	if err != nil {
		return err
	}

	err = p.profilePersistence.ChangePassword(ctx, changePasswordParam, userID)
	if err != nil {
		return err
	}

	p.logger.Info(ctx, "user changed password", zap.Any("user-id", userID))
	return nil
}

func (p *profileModule) GetAllCurrentSessions(ctx context.Context) ([]dto.InternalRefreshToken, error) {
	id, ok := ctx.Value(constant.Context("x-user-id")).(string)
	if !ok {
		err := errors.ErrInvalidUserInput.New("invalid user id")
		p.logger.Info(ctx, "invalid user id", zap.Error(err), zap.Any("user_id", id))
		return nil, err
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		err := errors.ErrNoRecordFound.Wrap(err, "user not found")
		p.logger.Info(ctx, "parse error", zap.Error(err), zap.String("user id", id))
		return nil, err
	}

	return p.oauthPersistence.GetInternalRefreshTokensByUserID(ctx, userID)
}
