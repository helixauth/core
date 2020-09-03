package oidc

import (
	"context"
	"os"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/app/oidc/api"
	"github.com/helixauth/helix/src/shared/gateway"

	"github.com/gin-gonic/gin"
)

func New(ctx context.Context, cfg config.Config, gateways gateway.Gateways) *gin.Engine {
	oidc := api.New(ctx, cfg, gateways)
	wd, _ := os.Getwd()
	public := wd + "/src/app/oidc/public"
	html := public + "/html/*"
	r := gin.New()
	r.Use(gin.Logger())
	r.LoadHTMLGlob(html)
	r.Static("/public", public)
	r.GET("/", oidc.Index)
	r.POST("/authenticate", oidc.Authenticate)
	r.GET("/authorization", oidc.Authorization)
	r.GET("/introspect", oidc.Introspect)
	r.GET("/jwks", oidc.JWKs)
	r.POST("/revoke", oidc.Revoke)
	r.POST("/token", oidc.Token)
	r.POST("/userinfo", oidc.UserInfo)
	r.GET("/.well-known/openid-configuration", oidc.Configuration)
	return r
}
