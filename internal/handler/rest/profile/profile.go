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
// @param image formData file  true "image"
// @Success      200  {object}  model.Response
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

// ChangePhone	 change's user phone number.
// @Summary      change user phone number.
// @Description  change user phone number.
// @Tags         profile
// @Accept       json
// @Produce      json
// @param ChangePhoneParam body dto.ChangePhoneParam true "ChangePhoneParam"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.ErrorResponse
// @Router       /profile/phone [patch]
// @Security	BearerAuth
func (p *profile) ChangePhone(ctx *gin.Context) {
	changePhoneParam := dto.ChangePhoneParam{}
	err := ctx.ShouldBind(&changePhoneParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.logger.Info(ctx, "unable to bind change phone information", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	requestCtx := ctx.Request.Context()

	err = p.profileModule.ChangePhone(requestCtx, changePhoneParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	p.logger.Info(ctx, "user changed phone", zap.Any("phone-to", changePhoneParam.Phone))
	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// ChangePassword	 change's user password.
// @Summary      change's user password..
// @Description  change's user password..
// @Tags         profile
// @Accept       json
// @Produce      json
// @param ChangePasswordParam body dto.ChangePasswordParam true "ChangePasswordParam"
// @Success      200  {object}  model.Response
// @Failure      400  {object}  model.ErrorResponse
// @Router       /profile/password [patch]
// @Security	BearerAuth
func (p *profile) ChangePassword(ctx *gin.Context) {
	changePasswordParam := dto.ChangePasswordParam{}
	err := ctx.ShouldBind(&changePasswordParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		p.logger.Info(ctx, "unable to bind change password information", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	requestCtx := ctx.Request.Context()

	err = p.profileModule.ChangePassword(requestCtx, changePasswordParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, nil, nil)
}

// GetAllCurrentSessions	 get's all current session of the user.
// @Summary     get's all current session of the user.
// @Description  get's all user sessions.
// @Tags         profile
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.InternalRefreshToken
// @Failure      400  {object}  model.ErrorResponse
// @Router       /profile/devices [get]
// @Security	BearerAuth
func (p *profile) GetAllCurrentSessions(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()

	sessions, err := p.profileModule.GetAllCurrentSessions(requestCtx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, sessions, nil)
}

func (p *profile) GetUserPermissions(ctx *gin.Context) {
	requestCtx := ctx.Request.Context()

	permissions, err := p.profileModule.GetUserPermissions(requestCtx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	constant.SuccessResponse(ctx, http.StatusOK, permissions, nil)
}
