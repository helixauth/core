package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/admin"
	"github.com/helixauth/helix/src/auth"
	"github.com/helixauth/helix/src/entity"
	"github.com/helixauth/helix/src/lib/database"
	"github.com/helixauth/helix/src/lib/email"

	"github.com/dchest/uniuri"
	_ "github.com/lib/pq"
)

func main() {
	ctx := cfg.Configure(context.Background())

	// secretsManager, err := secrets.New("cfg/secrets-dec.dev.yaml")
	// if err != nil {
	// 	panic(err)
	// }

	// Connect to database
	database, err := database.New(ctx)
	if err != nil {
		panic(err)
	}

	// Load tenant
	tenant := loadTenant(ctx, database)

	// Connect to email provider
	email, err := email.New(ctx, tenant)
	if err != nil {
		panic(err)
	}

	// Run apps
	go admin.Run(ctx, database)
	auth.Run(ctx, database, email)
}

func loadTenant(ctx context.Context, database database.Gateway) *entity.Tenant {

	// Get tenant ID
	tenantID, ok := ctx.Value(cfg.TenantID).(string)
	if !ok || tenantID == "" {
		panic("TENANT_ID not set")
	}

	// Check if tenant already exists
	log.Printf("🏠 Running as tenant: '%v'", tenantID)
	tenant := &entity.Tenant{}
	err := database.QueryItem(ctx, tenant, `SELECT * FROM tenants WHERE id = $1`, tenantID)
	if err == nil {
		return tenant
	} else if err != sql.ErrNoRows {
		panic(err)
	}

	// Create tenant
	txn, err := database.BeginTxn(ctx)
	if err != nil {
		panic(err)
	}
	tenant.ID = tenantID
	if err = txn.Insert(ctx, tenant); err != nil {
		panic(err)
	}

	// Create client
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

	return tenant
}
