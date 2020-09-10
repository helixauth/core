package app

import (
	"context"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/shared/database"
	"github.com/helixauth/helix/src/shared/email"

	"github.com/gin-gonic/gin"
)

type App interface {
	Authentication(c *gin.Context)
	Authorization(c *gin.Context)
	Configuration(c *gin.Context)
	Index(c *gin.Context)
	Introspect(c *gin.Context)
	JWKs(c *gin.Context)
	Revoke(c *gin.Context)
	Token(c *gin.Context)
	UserInfo(c *gin.Context)
}

type app struct {
	TenantID string
	Database database.Gateway
	Email    email.Gateway
}

func New(ctx context.Context, database database.Gateway, email email.Gateway) App {
	return &app{
		TenantID: ctx.Value(cfg.TenantID).(string),
		Database: database,
		Email:    email,
	}
}

func (a *app) context(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, cfg.TenantID, a.TenantID)
	return ctx
}
