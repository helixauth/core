package entity

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

// DatabaseTable points to the "tenants" table
func (t *Tenant) DatabaseTable() string {
	return "tenants"
}
