package oidc

import (
	"context"
	"os"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/oidc/app"
	"github.com/helixauth/helix/src/shared/gateway"

	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context, cfg config.Config, gateways gateway.Gateways) {
	app := app.New(ctx, cfg, gateways)
	wd, _ := os.Getwd()
	public := wd + "/src/oidc/public"
	html := public + "/html/*"
	r := gin.New()
	r.Use(gin.Logger())
	r.LoadHTMLGlob(html)
	r.Static("/public", public)
	r.GET("/", app.Index)
	r.POST("/authentication", app.Authentication)
	r.GET("/authorization", app.Authorization)
	r.GET("/introspect", app.Introspect)
	r.GET("/jwks", app.JWKs)
	r.POST("/revoke", app.Revoke)
	r.POST("/token", app.Token)
	r.POST("/userinfo", app.UserInfo)
	r.GET("/.well-known/openid-configuration", app.Configuration)
	r.Run(":80")
}
