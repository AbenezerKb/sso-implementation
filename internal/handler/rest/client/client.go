package client

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model/dto"
	"sso/internal/handler/rest"
	"sso/internal/module"
	"sso/platform/logger"
)

type client struct {
	logger       logger.Logger
	clientModule module.ClientModule
}

func Init(logger logger.Logger, clientModule module.ClientModule) rest.Client {
	return &client{
		logger:       logger,
		clientModule: clientModule,
	}
}

// CreateClient is a handler for creating a client
// @Summary      Create a client
// @Description  Create a new client
// @Tags         client
// @Accept       json
// @Produce      json
// @param client body dto.Client true "client"
// @Success      200  {object}  dto.Client
// @Failure      400  {object}  model.ErrorResponse
// @Router       /clients [post]
// @Security	BearerAuth
func (c *client) CreateClient(ctx *gin.Context) {
	clientParam := dto.Client{}
	err := ctx.ShouldBind(&clientParam)
	if err != nil {
		err := errors.ErrInvalidUserInput.Wrap(err, "invalid input")
		c.logger.Info(ctx, "couldn't bind to dto.Client body", zap.Error(err))
		_ = ctx.Error(err)
		return
	}

	createdClient, err := c.clientModule.Create(ctx.Request.Context(), clientParam)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	c.logger.Info(ctx, "created client")
	constant.SuccessResponse(ctx, http.StatusCreated, createdClient, nil)
}
