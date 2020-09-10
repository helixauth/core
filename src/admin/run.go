package admin

import (
	"context"
	"os"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/admin/app"
	"github.com/helixauth/helix/src/shared/gateway"

	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context, cfg config.Config, gateways gateway.Gateways) {
	app := app.New(ctx, cfg, gateways)
	wd, _ := os.Getwd()
	public := wd + "/src/admin/public"
	html := public + "/html/*"
	r := gin.New()
	r.Use(gin.Logger())
	r.LoadHTMLGlob(html)
	r.Static("/public", public)
	r.GET("/", app.Index)
	r.GET("/users", app.UsersList)
	r.GET("/users/:id", app.UsersGet)
	r.DELETE("/users/:id", app.UsersDelete)
	r.Run(":2048")
}
