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
	AWSRegion          = key("AWS_REGION")
	AWSAccessKeyID     = key("AWS_ACCESS_KEY_ID")
	AWSSecretAccessKey = key("AWS_SECRET_ACCESS_KEY")

	GinMode = key("GIN_MODE")

	PostgresUsername = key("POSTGRES_USERNAME")
	PostgresPassword = key("POSTGRES_PASSWORD")
	PostgresHost     = key("POSTGRES_HOST")
	PostgresPort     = key("POSTGRES_PORT")
	PostgresDBName   = key("POSTGRES_DB_NAME")
	PostgresSSLMode  = key("POSTGRES_SSL_MODE")

	TenantID = key("TENANT_ID")
)

var allKeys = []key{
	AWSRegion,
	AWSAccessKeyID,
	AWSSecretAccessKey,
	GinMode,
	PostgresUsername,
	PostgresPassword,
	PostgresHost,
	PostgresPort,
	PostgresDBName,
	PostgresSSLMode,
	TenantID,
}

func Configure(ctx context.Context) context.Context {
	for _, k := range allKeys {
		ctx = context.WithValue(ctx, k, os.Getenv(k.String()))
	}
	return ctx
}
