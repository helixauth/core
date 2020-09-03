package main

import (
	"context"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/app/admin"
	"github.com/helixauth/helix/src/app/oidc"
	"github.com/helixauth/helix/src/shared/gateway"
)

func main() {
	ctx := context.Background()
	cfg := config.New()
	gateways, err := gateway.New(cfg)
	if err != nil {
		panic(err)
	}
	go admin.New(ctx, cfg, gateways).Run(":2048")
	oidc.New(ctx, cfg, gateways).Run(":80")
}
