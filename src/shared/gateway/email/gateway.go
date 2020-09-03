package email

import (
	"context"

	"github.com/helixauth/helix/cfg"
)

const (
	charSet string = "UTF-8"
)

type Gateway interface {
	SendEmail(ctx context.Context, sender string, recipient string, subject string, htmlBody string) error
}

func New(cfg config.Config) (Gateway, error) {
	if cfg.Email.SES != nil {
		return NewSESGateway(cfg)
	}
	return NewFakeGateway(cfg)
}
