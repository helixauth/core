package email

import (
	"context"

	"github.com/helixauth/helix/src/entity"
	"github.com/helixauth/helix/src/lib/mapper"
)

const (
	charSet string = "UTF-8"
)

const (
	SES = "SES"
)

type Gateway interface {
	SendEmail(ctx context.Context, sender string, recipient string, subject string, htmlBody string) error
}

func New(ctx context.Context, tenant *entity.Tenant) (Gateway, error) {
	switch mapper.String(tenant.EmailProvider) {
	case SES:
		return NewSESGateway(ctx, tenant)

	default:
		return NewFakeGateway(ctx)

	}
}
