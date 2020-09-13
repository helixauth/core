package auth

import (
	"context"
	"html/template"
	"os"

	"github.com/helixauth/helix/src/auth/app"
	"github.com/helixauth/helix/src/lib/database"
	"github.com/helixauth/helix/src/lib/email"

	"github.com/gin-gonic/gin"
)

// Run starts the auth application (OAuth2/OIDC server)
func Run(ctx context.Context, database database.Gateway, email email.Gateway) {
	app := app.New(ctx, database, email)
	wd, _ := os.Getwd()
	public := wd + "/src/auth/public"
	html := public + "/html/*"
	r := gin.New()
	r.Use(gin.Logger())
	r.SetFuncMap(template.FuncMap{
		"safeURL": func(u string) template.URL { return template.URL(u) },
	})
	r.LoadHTMLGlob(html)
	r.Static("/public", public)
	r.GET("/", app.Index)
	r.GET("/authorize", app.Authorize)
	r.POST("/authorize", app.Authorize)
	r.GET("/introspect", app.Introspect)
	r.GET("/jwks", app.JWKs)
	r.POST("/revoke", app.Revoke)
	r.POST("/token", app.Token)
	r.POST("/userinfo", app.UserInfo)
	r.GET("/.well-known/openid-configuration", app.Configuration)
	r.Run(":80")
}
