package entity

// Client represents an application that has delegated auth responsibilities to a tenant
type Client struct {
	ID                string  `json:"id"`
	TenantID          string  `json:"tenant_id"`
	Name              *string `json:"name"`
	Secret            *string `json:"secret"`
	Picture           *string `json:"picture"`
	Website           *string `json:"website"`
	Description       *string `json:"description"`
	PrivacyPolicy     *string `json:"privacy_policy"`
	IsThirdParty      bool    `json:"is_third_party"`
	AuthorizedDomains string  `json:"authorized_domains"`
}

// DatabaseTable points to the "clients" table
func (c *Client) DatabaseTable() string {
	return "clients"
}
