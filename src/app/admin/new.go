package admin

import (
	"context"
	"os"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/app/admin/api"
	"github.com/helixauth/helix/src/shared/gateway"

	"github.com/gin-gonic/gin"
)

func New(ctx context.Context, cfg config.Config, gateways gateway.Gateways) *gin.Engine {
	admin := api.New(ctx, cfg, gateways)
	wd, _ := os.Getwd()
	public := wd + "/src/app/admin/public"
	html := public + "/html/*"
	r := gin.New()
	r.Use(gin.Logger())
	r.LoadHTMLGlob(html)
	r.Static("/public", public)
	r.GET("/", admin.Index)
	r.GET("/users", admin.UsersList)
	r.GET("/users/:id", admin.UsersGet)
	r.DELETE("/users/:id", admin.UsersDelete)
	return r
}
