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
)

func main() {
	ctx := cfg.Configure(context.Background())

	// Connect to the database
	database, err := database.New(ctx)
	if err != nil {
		panic(err)
	}

	// Load the tenant
	loadTenant(ctx, database)

	// TODO make this dynamic
	email, err := email.New(ctx)
	if err != nil {
		panic(err)
	}

	// Run apps
	go admin.Run(ctx, database)
	oidc.Run(ctx, database, email)
}

<<<<<<< HEAD
func bootstrap(ctx context.Context, database database.Gateway) {

	// Check for an existing tenant
=======
func loadTenant(ctx context.Context, database database.Gateway) {
	// Get tenant ID
>>>>>>> 508d240652c0491e9991f16b4570d36ca3c565c6
	tenantID, ok := ctx.Value(cfg.TenantID).(string)
	if !ok || tenantID == "" {
		panic("TENANT_ID not set")
	}
<<<<<<< HEAD
=======

	// Check if tenant already exists
>>>>>>> 508d240652c0491e9991f16b4570d36ca3c565c6
	log.Printf("🏠 Running as tenant: '%v'", tenantID)
	tenant := &entity.Tenant{}
	if err := database.Query(ctx, tenant, `SELECT * FROM tenants WHERE id = $1`, tenantID); err != nil {
		panic(err)
	} else if (*tenant != entity.Tenant{}) {
		return
	}

<<<<<<< HEAD
	// Create a new tenant
	tx, err := database.BeginTx(ctx)
=======
	// Create tenant
	txn, err := database.Txn(ctx)
>>>>>>> 508d240652c0491e9991f16b4570d36ca3c565c6
	if err != nil {
		panic(err)
	}
	tenant.ID = tenantID
	if err = txn.Insert(ctx, tenant); err != nil {
		txn.Rollback()
		panic(err)
	}
<<<<<<< HEAD
	if err = tx.Commit(); err != nil {
=======
	if err = txn.Commit(); err != nil {
>>>>>>> 508d240652c0491e9991f16b4570d36ca3c565c6
		panic(err)
	}
}
