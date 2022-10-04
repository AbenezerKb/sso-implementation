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

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type profileModule struct {
	logger                logger.Logger
	oauthPersistence      storage.OAuthPersistence
	profilePersistence    storage.ProfilePersistence
	profilePictureDist    string
	profilePictureMaxSize int
}

func InitProfile(logger logger.Logger, oauthPersistence storage.OAuthPersistence, profilePersistence storage.ProfilePersistence, profilePictureDist string, profilePictureMaxSize int) module.ProfileModule {
	return &profileModule{
		logger:                logger,
		oauthPersistence:      oauthPersistence,
		profilePersistence:    profilePersistence,
		profilePictureDist:    profilePictureDist,
		profilePictureMaxSize: profilePictureMaxSize,
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

	if imageFile.Size > int64(p.profilePictureMaxSize) {
		err = errors.ErrInvalidUserInput.New("image size must be less than 2MB")
		p.logger.Info(ctx, "image size too big", zap.Error(err), zap.String("image", imageFile.Filename), zap.Any("size", imageFile.Size), zap.Any("user-id", id))
		return err
	}

	// final image name
	finalImageName := fmt.Sprint(time.Now().UnixMilli()) + "_" + id + "_" + imageFile.Filename

	err = utils.SaveMultiPartFile(imageFile, p.profilePictureDist+finalImageName)
	if err != nil {
		err = errors.ErrInternalServerError.Wrap(err, "couldn't save profile picture")
		p.logger.Error(context.Background(), "error unable to save profile picture to disck", zap.Error(err), zap.Any("image", imageFile))
	}

	err = p.profilePersistence.UpdateProfilePicture(ctx, finalImageName, userID)
	if err != nil {
		return err
	}

	return nil
}
