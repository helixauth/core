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

	"github.com/dchest/uniuri"
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

func loadTenant(ctx context.Context, database database.Gateway) {
	// Get tenant ID
	tenantID, ok := ctx.Value(cfg.TenantID).(string)
	if !ok || tenantID == "" {
		panic("TENANT_ID not set")
	}

	// Check if tenant already exists
	log.Printf("🏠 Running as tenant: '%v'", tenantID)
	tenant := &entity.Tenant{}
	if err := database.Query(ctx, tenant, `SELECT * FROM tenants WHERE id = $1`, tenantID); err != nil {
		panic(err)
	} else if (*tenant != entity.Tenant{}) {
		return
	}

	// Create tenant
	txn, err := database.Txn(ctx)
	if err != nil {
		panic(err)
	}
	tenant.ID = tenantID
	if err = txn.Insert(ctx, tenant); err != nil {
		panic(err)
	}

	// Create first-party client
	client := &entity.Client{
		ID:                uniuri.NewLen(32),
		TenantID:          tenantID,
		Name:              nil,
		Secret:            nil,
		Picture:           nil,
		Website:           nil,
		Description:       nil,
		PrivacyPolicy:     nil,
		IsThirdParty:      false,
		AuthorizedDomains: []string{"localhost:3000"},
	}
	if err = txn.Insert(ctx, client); err != nil {
		panic(err)
	}

	if err = txn.Commit(); err != nil {
		panic(err)
	}
}
