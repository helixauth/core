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
	TenantID         = key("TENANT_ID")
	PostgresUsername = key("POSTGRES_USERNAME")
	PostgresPassword = key("POSTGRES_PASSWORD")
	PostgresHost     = key("POSTGRES_HOST")
	PostgresPort     = key("POSTGRES_PORT")
	PostgresDBName   = key("POSTGRES_DB_NAME")
	PostgresSSLMode  = key("POSTGRES_SSL_MODE")
)

var allKeys = []key{
	TenantID,
	PostgresUsername,
	PostgresPassword,
	PostgresHost,
	PostgresPort,
	PostgresDBName,
	PostgresSSLMode,
}

func Configure(ctx context.Context) context.Context {
	for _, k := range allKeys {
		ctx = context.WithValue(ctx, k, os.Getenv(k.String()))
	}
	return ctx
}
