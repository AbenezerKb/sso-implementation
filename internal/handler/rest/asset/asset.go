package asset

import (
	"net/http"

	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type asset struct {
	log   logger.Logger
	asset module.Asset
}

func Init(log logger.Logger, assetModule module.Asset) rest.Asset {
	return &asset{
		log:   log,
		asset: assetModule,
	}
}

// UploadAsset uploads an asset file to the server storage
// @Summary      uploads an asset file to the server storage
// @Description  uploads an asset file to the server storage
// @Tags         asset
// @Accept       mpfd
// @Produce      json
// @param type body string true "type"
// @param asset body string true "asset"
// @Success      201  {string} string
// @Router 		 /assets [post]
func (a *asset) UploadAsset(ctx *gin.Context) {
	var uploadRequest dto.UploadAssetRequest

	err := ctx.ShouldBind(&uploadRequest)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		a.log.Info(ctx, "failed to bind asset upload request", zap.Error(err))
		_ = ctx.Error(err)

		return
	}

	fileName, err := a.asset.UploadAsset(ctx, uploadRequest)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusCreated, fileName, nil)
}
