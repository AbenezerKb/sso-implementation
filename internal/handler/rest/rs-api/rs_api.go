package rs_api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto/request_models"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"
)

type rsAPI struct {
	logger logger.Logger
	rsAPI  module.RSAPI
}

func Init(logger logger.Logger, rs module.RSAPI) rest.RSAPI {
	return &rsAPI{
		logger,
		rs,
	}
}

// GetUserByPhoneOrID	 returns a user with the specified id or phone.
// @Summary      returns a user with id or phone number
// @Description  returns a user with id or phone number.
// @Tags         internal
// @Accept       json
// @Produce      json
// @param id query string true "id"
// @param phone query string true "phone"
// @Success      200  {object}  dto.User
// @Failure      400  {object}  model.ErrorResponse
// @Router       /internal/user [get]
// @Security	BasicAuth
func (i *rsAPI) GetUserByPhoneOrID(ctx *gin.Context) {
	var req request_models.RSAPIUserRequest
	err := ctx.BindQuery(&req)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid query params")
		i.logger.Info(ctx, "invalid query params for rs-api")
		_ = ctx.Error(err)
		return
	}

	requestCtx := ctx.Request.Context()
	user, err := i.rsAPI.GetUserByIDOrPhone(requestCtx, req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	i.logger.Info(ctx, "user details fetched", zap.Any("request", req))
	constant.SuccessResponse(ctx, http.StatusOK, user, nil)
}
