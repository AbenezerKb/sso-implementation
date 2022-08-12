package user

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

type user struct {
	logger     logger.Logger
	userModule module.UserModule
}

func Init(logger logger.Logger, userModule module.UserModule) rest.User {
	return &user{
		logger,
		userModule,
	}
}

// CreateUser	 creates a new user.
// @Summary      create a new user.
// @Description  create a new user.
// @Tags         user
// @Accept       json
// @Produce      json
// @param user body dto.CreateUser true "user"
// @Success      200  {object}  dto.User
// @Failure      400  {object}  model.ErrorResponse
// @Router       /users [post]
func (u *user) CreateUser(ctx *gin.Context) {
	userParam := dto.CreateUser{}
	err := ctx.ShouldBind(&userParam)
	if err != nil {
		u.logger.Info(ctx, "unable to bind user data", zap.Error(err))
		_ = ctx.Error(errors.ErrInvalidUserInput.Wrap(err, "invalid input"))
		return
	}
	createdUser, err := u.userModule.Create(ctx.Request.Context(), userParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}
	u.logger.Info(ctx, "created user")
	constant.SuccessResponse(ctx, http.StatusCreated, createdUser, nil)
}
