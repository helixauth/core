package app

import (
	"context"

	"github.com/helixauth/core/cfg"
	"github.com/helixauth/core/src/lib/database"
	"github.com/helixauth/core/src/lib/email"
	"github.com/helixauth/core/src/lib/secrets"

	"github.com/gin-gonic/gin"
)

// App provides an interface to the auth application (i.e. OAuth2/OIDC server)
type App interface {
	Authorize(c *gin.Context)
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
	Secrets  secrets.Manager
}

// New creates a new auth application
func New(ctx context.Context,
	database database.Gateway,
	email email.Gateway,
	secrets secrets.Manager) App {
	return &app{
		TenantID: ctx.Value(cfg.TenantID).(string),
		Database: database,
		Email:    email,
		Secrets:  secrets,
	}
}

func (a *app) context(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, cfg.TenantID, a.TenantID)
	return ctx
}
