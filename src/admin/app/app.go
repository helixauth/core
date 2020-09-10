package app

import (
	"context"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/shared/gateway"

	"github.com/gin-gonic/gin"
)

type App interface {
	Index(c *gin.Context)

	UsersDelete(c *gin.Context)
	UsersGet(c *gin.Context)
	UsersList(c *gin.Context)
}

type app struct {
	Config   config.Config
	Gateways gateway.Gateways
}

func New(ctx context.Context, cfg config.Config, gateways gateway.Gateways) App {
	return &app{
		Config:   cfg,
		Gateways: gateways,
	}
}
