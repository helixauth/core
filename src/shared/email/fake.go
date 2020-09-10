package email

import (
	"context"
	"log"
)

type fakeGateway struct{}

func NewFakeGateway(ctx context.Context) (Gateway, error) {
	return &fakeGateway{}, nil
}

func (g *fakeGateway) SendEmail(ctx context.Context, sender string, recipient string, subject string, htmlBody string) error {
	log.Printf(`
    *** FAKE EMAIL GATEWAY ***

    from: %v
    to: %v
    subject: %v
    body: %v

  `, sender, recipient, subject, htmlBody)
	return nil
}
