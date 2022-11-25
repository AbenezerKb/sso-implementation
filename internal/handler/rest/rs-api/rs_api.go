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
// @param user query request_models.RSAPIUserRequest true "user"
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

// GetUsersByPhoneOrID	 returns users with the specified ids or phones.
// @Summary      returns users with ids or phone numbers
// @Description  returns user with ids or phone numbers.
// @Tags         internal
// @Accept       json
// @Produce      json
// @param users body request_models.RSAPIUsersRequest true "users"
// @Success      200  {object}  []dto.User
// @Failure      400  {object}  model.ErrorResponse
// @Router       /internal/users [get]
// @Security	BasicAuth
func (i *rsAPI) GetUsersByPhoneOrID(ctx *gin.Context) {
	var req request_models.RSAPIUsersRequest
	err := ctx.ShouldBind(&req)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid request body")
		i.logger.Info(ctx, "invalid request body for rs-api")
		_ = ctx.Error(err)
		return
	}

	requestCtx := ctx.Request.Context()
	user, err := i.rsAPI.GetUsersByIDOrPhone(requestCtx, req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	i.logger.Info(ctx, "users detail fetched")
	constant.SuccessResponse(ctx, http.StatusOK, user, nil)
}
