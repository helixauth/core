package main

import (
	"context"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/admin"
	"github.com/helixauth/helix/src/oidc"
	"github.com/helixauth/helix/src/shared/database"
	"github.com/helixauth/helix/src/shared/email"
)

func main() {
	ctx := cfg.Configure(context.Background())

	// Build dependencies
	database, err := database.New(ctx)
	if err != nil {
		panic(err)
	}
	email, err := email.New(ctx)
	if err != nil {
		panic(err)
	}

	// Run apps
	go admin.Run(ctx, database)
	oidc.Run(ctx, database, email)
}
