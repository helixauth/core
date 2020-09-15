package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/helixauth/core/cfg"
	"github.com/helixauth/core/src/admin"
	"github.com/helixauth/core/src/auth"
	"github.com/helixauth/core/src/entity"
	"github.com/helixauth/core/src/lib/database"
	"github.com/helixauth/core/src/lib/email"
	"github.com/helixauth/core/src/lib/secrets"

	"github.com/dchest/uniuri"
	_ "github.com/lib/pq"
)

func main() {
	ctx := cfg.Configure(context.Background())

	// Load secrets
	// TODO pick secrets file based on HELIX env
	secrets, err := secrets.New("cfg/secrets.dec.dev.yaml")
	if err != nil {
		panic(err)
	}

	// Connect to database
	connInfo, err := getDatabaseConnInfo(secrets)
	if err != nil {
		panic(err)
	}
	database, err := database.New(ctx, connInfo)
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
	auth.Run(ctx, database, email, secrets)
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

func getDatabaseConnInfo(secrets secrets.Manager) (database.ConnInfo, error) {
	var err error
	connInfo := database.ConnInfo{}
	connArgs := map[string]interface{}{
		"postgres.host":     "",
		"postgres.port":     "",
		"postgres.username": "",
		"postgres.password": "",
		"postgres.db_name":  "",
		"postgres.ssl_mode": "",
	}
	for k := range connArgs {
		if connArgs[k], err = secrets.Get(k); err != nil {
			return connInfo, err
		}
	}
	connInfo.Host = connArgs["postgres.host"].(string)
	connInfo.Port = connArgs["postgres.port"].(int)
	connInfo.Username = connArgs["postgres.username"].(string)
	connInfo.Password = connArgs["postgres.password"].(string)
	connInfo.DBName = connArgs["postgres.db_name"].(string)
	connInfo.SSLMode = connArgs["postgres.ssl_mode"].(string)
	return connInfo, nil
}
