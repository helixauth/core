package api

import (
	"context"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/shared/gateway"

	"github.com/gin-gonic/gin"
)

type API interface {
	Authenticate(c *gin.Context)
	Authorization(c *gin.Context)
	Configuration(c *gin.Context)
	Index(c *gin.Context)
	Introspect(c *gin.Context)
	JWKs(c *gin.Context)
	Revoke(c *gin.Context)
	Token(c *gin.Context)
	UserInfo(c *gin.Context)
}

type api struct {
	Config   config.Config
	Gateways gateway.Gateways
}

func New(ctx context.Context, cfg config.Config, gateways gateway.Gateways) API {
	return &api{
		Config:   cfg,
		Gateways: gateways,
	}
}
