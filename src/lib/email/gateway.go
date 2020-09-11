package email

import (
	"context"
)

const (
	charSet string = "UTF-8"
)

type Gateway interface {
	SendEmail(ctx context.Context, sender string, recipient string, subject string, htmlBody string) error
}

func New(ctx context.Context) (Gateway, error) {

	// TODO load SES gateway

	// if cfg.Email.SES != nil {
	// 	return NewSESGateway(cfg)
	// }

	return NewFakeGateway(ctx)
}
