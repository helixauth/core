package entity

import (
	"database/sql"

	"github.com/helixauth/helix/src/shared/utils"
)

// Tenant represents an application which can issue user identities
type Tenant struct {
	ID                 string  `json:"id"`
	Name               *string `json:"name"`
	Picture            *string `json:"picture"`
	Website            *string `json:"website"`
	Email              *string `json:"email"`
	EmailProvider      *string `json:"email_provider"`
	AWSRegionID        *string `json:"aws_region_id"`
	AWSAccessKeyID     *string `json:"aws_access_key_id"`
	AWSSecretAccessKey *string `json:"aws_secret_access_key"`
}

// FromSQL parses a Tenant entity from a SQL row
func (t *Tenant) FromSQL(rows *sql.Rows) error {
	if !rows.Next() {
		return nil
	}
	return utils.SQLParseRow(rows, t)
}

// SQLTable points to the "tenants" table
func (t *Tenant) SQLTable() string {
	return "tenants"
}
