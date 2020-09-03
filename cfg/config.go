package config

import (
	"os"
	"strings"
)

func New() Config {
	return Config{
		Email: EmailConfig{
			SES: &SESEmailConfig{
				Region:          os.Getenv("AWS_REGION"),
				AcessKeyID:      os.Getenv("AWS_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			},
		},
		Gin: GinConfig{
			Mode: os.Getenv("GIN_MODE"),
		},
		OAuth: OAuthConfig{
			AuthorizedDomains: strings.Split(os.Getenv("OAUTH_AUTHORIZED_DOMAINS"), " "),
			ClientID:          os.Getenv("OAUTH_CLIENT_ID"),
			ClientSecret:      os.Getenv("OAUTH_CLIENT_SECRET"),
		},
		Postgres: PostgresConfig{
			Username: os.Getenv("POSTGRES_USERNAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			DBName:   os.Getenv("POSTGRES_DB_NAME"),
			SSLMode:  os.Getenv("POSTGRES_SSL_MODE"),
		},
	}
}

type Config struct {
	Email    EmailConfig
	Gin      GinConfig
	OAuth    OAuthConfig
	Postgres PostgresConfig
}

type EmailConfig struct {
	SES *SESEmailConfig
}

type SESEmailConfig struct {
	Region          string
	AcessKeyID      string
	SecretAccessKey string
}

type GinConfig struct {
	Mode string
}

type OAuthConfig struct {
	AuthorizedDomains []string
	ClientID          string
	ClientSecret      string
}

type PostgresConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}
