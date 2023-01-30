package asset

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/constant/state"
	"sso/internal/module"
	"sso/platform"
	"sso/platform/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm/utils"
)

type asset struct {
	log         logger.Logger
	fileManager platform.Asset
	options     state.UploadParams
}

func Init(log logger.Logger, fileManager platform.Asset, options state.UploadParams) module.Asset {
	return &asset{
		log:         log,
		fileManager: fileManager,
		options:     options,
	}
}

func (a *asset) UploadAsset(ctx context.Context, param dto.UploadAssetRequest) (string, error) {
	if err := param.Validate(); err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		a.log.Info(ctx, "invalid input for upload asset", zap.Error(err))
		return "", err
	}

	var fileType *state.FileType
	for i := 0; i < len(a.options.FileTypes); i++ {
		if param.Type == a.options.FileTypes[i].Name { // if it is one of the specified types
			fileType = &a.options.FileTypes[i]

			break
		}
	}
	if fileType == nil {
		err := errors.ErrInvalidUserInput.New(fmt.Sprintf("type %s does not exist", param.Type))
		a.log.Info(ctx, "invalid upload with non-existing upload type", zap.String("type", param.Type))
		return "", err
	}

	src, err := param.Asset.Open()
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "couldn't open asset")
		a.log.Warn(ctx, "invalid file input for uploading asset", zap.Error(err))
		return "", err
	}

	if param.Asset.Size > fileType.MaxSize {
		err := errors.ErrInvalidUserInput.New(
			fmt.Sprintf("%s size must be less than %dMB",
				fileType.Name, fileType.MaxSize/1024/1024))
		a.log.Info(ctx, "asset upload greater than allowed size", zap.Error(err))
		return "", err
	}

	temp := strings.Split(param.Asset.Filename, ".")
	if len(temp) < 2 {
		err := errors.ErrInvalidUserInput.New("asset has no file type")
		a.log.Info(ctx, "asset upload with no extension", zap.Error(err))
		return "", err
	}

	fileExtension := temp[len(temp)-1]
	if !utils.Contains(fileType.Types, fileExtension) {
		err := errors.ErrInvalidUserInput.New(
			fmt.Sprintf("%s must be one of types (%s)", fileType.Name, strings.Join(fileType.Types, ",")))
		a.log.Error(ctx, "asset upload with invalid file type", zap.Error(err))
		return "", err
	}

	newFileName := fmt.Sprintf("%s-%d.%s",
		uuid.NewString(),
		time.Now().Unix(),
		fileExtension,
	)
	err = a.fileManager.SaveAsset(ctx, src, path.Join(fileType.Name, newFileName))
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "unable to save asset")
		a.log.Error(ctx, "error while saving asset file", zap.Error(err))
		return "", err
	}

	return newFileName, nil
}
