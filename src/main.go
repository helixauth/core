package main

import (
	"context"
	"log"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/admin"
	"github.com/helixauth/helix/src/oidc"
	"github.com/helixauth/helix/src/shared/database"
	"github.com/helixauth/helix/src/shared/email"
	"github.com/helixauth/helix/src/shared/entity"
	"github.com/helixauth/helix/src/shared/utils"
)

func main() {
	ctx := cfg.Configure(context.Background())

	// Connect to the database
	database, err := database.New(ctx)
	if err != nil {
		panic(err)
	}

	// Bootstrap the tenant
	bootstrap(ctx, database)

	// TODO make this dynamic
	email, err := email.New(ctx)
	if err != nil {
		panic(err)
	}

	// Run apps
	go admin.Run(ctx, database)
	oidc.Run(ctx, database, email)
}

func bootstrap(ctx context.Context, database database.Gateway) {

	// Check for an existing tenant
	tenantID, ok := ctx.Value(cfg.TenantID).(string)
	if !ok || tenantID == "" {
		panic("TENANT_ID not set")
	}
	log.Printf("🏠 Running as tenant: '%v'", tenantID)
	tenant := &entity.Tenant{}
	if err := database.Query(ctx, tenant, `SELECT * FROM tenants WHERE id = $1`, tenantID); err != nil {
		panic(err)
	} else if (*tenant != entity.Tenant{}) {
		return
	}

	// Create a new tenant
	tx, err := database.BeginTx(ctx)
	if err != nil {
		panic(err)
	}
	tenant.ID = tenantID
	if err = utils.SQLInsert(ctx, tenant, "tenants", tx); err != nil {
		tx.Rollback()
		panic(err)
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}
}
