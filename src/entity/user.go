package entity

import (
	"time"
)

// User represents a person using a client application
type User struct {
	ID                string            `json:"id"`
	TenantID          string            `json:"tenant_id"`
	Name              *string           `json:"name"`
	Nickname          *string           `json:"nickname"`
	PreferredUsername *string           `json:"preferred_username"`
	GivenName         *string           `json:"given_name"`
	MiddleName        *string           `json:"middle_name"`
	FamilyName        *string           `json:"family_name"`
	Email             *string           `json:"email"`
	EmailVerified     *bool             `json:"email_verified"`
	ZoneInfo          *string           `json:"zone_info"`
	Locale            *string           `json:"locale"`
	Address           *string           `json:"address"`
	PhoneNumber       *string           `json:"phone_number"`
	Picture           *string           `json:"picture"`
	Website           *string           `json:"website"`
	Gender            *string           `json:"gender"`
	Birthdate         *string           `json:"birthdate"`
	PasswordHash      *string           `json:"password_hash"`
	IsBlocked         bool              `json:"is_blocked"`
	Metadata          map[string]string `json:"metadata"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	LastActiveAt      *time.Time        `json:"last_active_at"`
}

// DatabaseTable points to the "users" table
func (u *User) DatabaseTable() string {
	return "users"
}
