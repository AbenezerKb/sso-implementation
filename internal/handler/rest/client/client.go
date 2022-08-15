package client

import (
	"github.com/gin-gonic/gin"
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

}
