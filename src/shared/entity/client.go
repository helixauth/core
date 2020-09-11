package entity

import (
	"database/sql"

	"github.com/helixauth/helix/src/shared/utils"
)

type Client struct {
	ID                string   `json:"id"`
	TenantID          string   `json:"tenant_id"`
	Name              *string  `json:"name"`
	Secret            *string  `json:"secret"`
	Picture           *string  `json:"picture"`
	Website           *string  `json:"website"`
	Description       *string  `json:"description"`
	PrivacyPolicy     *string  `json:"privacy_policy"`
	IsThirdParty      bool     `json:"is_third_party"`
	AuthorizedDomains []string `json:"authorized_domains"`
}

// FromSQL parses a Client entity from a SQL row
func (c *Client) FromSQL(rows *sql.Rows) error {
	if !rows.Next() {
		return nil
	}
	return utils.SQLParseRow(rows, c)
}

func (c *Client) SQLTable() string {
	return "clients"
}
