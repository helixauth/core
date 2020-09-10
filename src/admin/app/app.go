package app

import (
	"context"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/shared/database"

	"github.com/gin-gonic/gin"
)

type App interface {
	Index(c *gin.Context)

	UsersDelete(c *gin.Context)
	UsersGet(c *gin.Context)
	UsersList(c *gin.Context)
}

type app struct {
	TenantID string
	Database database.Gateway
}

func New(ctx context.Context, database database.Gateway) App {
	return &app{
		TenantID: ctx.Value(cfg.TenantID).(string),
		Database: database,
	}
}

func (a *app) context(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, cfg.TenantID, a.TenantID)
	return ctx
}
