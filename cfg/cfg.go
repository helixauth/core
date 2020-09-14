package cfg

import (
	"context"
	"os"
)

type key string

func (k key) String() string {
	return string(k)
}

const (
	TenantID = key("TENANT_ID")
)

var allKeys = []key{
	TenantID,
}

func Configure(ctx context.Context) context.Context {
	for _, k := range allKeys {
		ctx = context.WithValue(ctx, k, os.Getenv(k.String()))
	}
	return ctx
}
