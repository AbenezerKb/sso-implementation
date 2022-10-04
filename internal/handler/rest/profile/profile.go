package profile

import (
	"context"
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
// @Router       /profile [put]
// @Security	BearerAuth
func (p *profile) UpdateProfile(ctx *gin.Context) {
	userParam := dto.User{}
	err := ctx.ShouldBind(&userParam)
	if err != nil {
		p.logger.Info(ctx, "unable to bind user data", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}
	requestCtx := ctx.Request.Context()

	updatedUser, err := p.profileModule.UpdateProfile(requestCtx, userParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	p.logger.Info(ctx, "user profile updated", zap.Any("user", userParam))
	constant.SuccessResponse(ctx, http.StatusOK, updatedUser, nil)
}

// GetProfile	 get's user's profile.
// @Summary      get's user profile.
// @Description  get's user profile.
// @Tags         profile
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.User
// @Failure      400  {object}  model.ErrorResponse
// @Router       /profile [get]
// @Security	BearerAuth
func (p *profile) GetProfile(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()

	user, err := p.profileModule.GetProfile(requestCtx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	p.logger.Info(ctx, "user profile fetched", zap.Any("user", user))
	constant.SuccessResponse(ctx, http.StatusOK, user, nil)
}

// UpdateProfilePicture	 updates user's profile picture.
// @Summary      update user profile picture.
// @Description  update user profile picture.
// @Tags         profile
// @Accept       mpfd
// @Produce      json
// @param image formData  true "image"
// @Success      200  {object}
// @Failure      400  {object}  model.ErrorResponse
// @Router       /profile/picture [put]
// @Security	BearerAuth
func (p *profile) UpdateProfilePicture(ctx *gin.Context) {
	imageFile, err := ctx.FormFile("image")
	if err != nil {
		err = errors.ErrInvalidUserInput.Wrap(err, "invalid profile picture")
		p.logger.Error(context.Background(), "error binding profile picture")
		_ = ctx.Error(err)
		return
	}

	requestCtx := ctx.Request.Context()
	err = p.profileModule.UpdateProfilePicture(requestCtx, imageFile)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	p.logger.Info(ctx, "user profile picture updated", zap.Any("user-id", constant.Context("x-user-id")), zap.Any("picture", imageFile.Filename))
	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}
