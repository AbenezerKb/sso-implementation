package profile

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

type profile struct {
	logger        logger.Logger
	profileModule module.ProfileModule
}

func Init(logger logger.Logger, profileModule module.ProfileModule) rest.Profile {
	return &profile{
		logger:        logger,
		profileModule: profileModule,
	}
}

// UpdateProfile	 updates user's profile.
// @Summary      update user profile.
// @Description  update user profile.
// @Tags         profile
// @Accept       json
// @Produce      json
// @param user body dto.User true "user"
// @Success      200  {object}  dto.User
// @Failure      400  {object}  model.ErrorResponse
// @Router       /users [put]
// @Security	BearerAuth
func (u *profile) UpdateProfile(ctx *gin.Context) {
	userParam := dto.User{}
	err := ctx.ShouldBind(&userParam)
	if err != nil {
		u.logger.Info(ctx, "unable to bind user data", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}
	requestCtx := ctx.Request.Context()

	updatedUser, err := u.profileModule.UpdateProfile(requestCtx, userParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	u.logger.Info(ctx, "user profile updated", zap.Any("user", userParam))
	constant.SuccessResponse(ctx, http.StatusOK, updatedUser, nil)
}
