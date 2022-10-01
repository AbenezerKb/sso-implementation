package mini_ride

import (
	"net/http"
	"sso/internal/constant"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type miniRide struct {
	logger         logger.Logger
	miniRideModule module.MiniRideModule
}

func Init(logger logger.Logger, miniRideModule module.MiniRideModule) rest.MiniRide {
	return &miniRide{
		logger,
		miniRideModule,
	}
}

// CheckPhone	 check's if phone exists.
// @Summary      check's if phone exists.
// @Description  check's if phone exists.
// @Tags         miniRide
// @Accept       json
// @Produce      json
// @param phone path string true "phone"
// @Success      200  {object}  dto.MiniRideResponse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Router       /users/{phone}/exists [get]
// @Security	BasicAuth
func (m *miniRide) CheckPhone(ctx *gin.Context) {
	phone := ctx.Param("phone")

	requestCtx := ctx.Request.Context()
	rsp, err := m.miniRideModule.CheckPhone(requestCtx, phone)

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	m.logger.Info(ctx, "user checked by mini-ride", zap.Any("phone", phone))
	constant.SuccessResponse(ctx, http.StatusOK, rsp, nil)
}
