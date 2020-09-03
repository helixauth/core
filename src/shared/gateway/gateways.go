package gateway

import (
	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/shared/gateway/database"
	"github.com/helixauth/helix/src/shared/gateway/email"
)

type Gateways struct {
	Database database.Gateway
	Email    email.Gateway
}

func New(cfg config.Config) (Gateways, error) {
	databaseGateway, err := database.New(cfg)
	if err != nil {
		return Gateways{}, err
	}
	emailGateway, err := email.New(cfg)
	if err != nil {
		return Gateways{}, err
	}
	return Gateways{
		Database: databaseGateway,
		Email:    emailGateway,
	}, nil
}
