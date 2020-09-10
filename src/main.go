package main

import (
	"context"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/admin"
	"github.com/helixauth/helix/src/oidc"
	"github.com/helixauth/helix/src/shared/gateway"
)

func main() {
	ctx := context.Background()
	cfg := config.New()
	gateways, err := gateway.New(cfg)
	if err != nil {
		panic(err)
	}

	go admin.Run(ctx, cfg, gateways)
	oidc.Run(ctx, cfg, gateways)
}
