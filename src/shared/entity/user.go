package entity

import (
	"database/sql"
	"time"

	"github.com/helixauth/helix/src/shared/utils"
)

// TODO add TenantID

type User struct {
	ID                string            `json:"id"`
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
	IsBlocked         bool              `json:"is_blocked"`
	Metadata          map[string]string `json:"metadata"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	LastActiveAt      *time.Time        `json:"last_active_at"`
}

func (u *User) FromSQL(rows *sql.Rows) error {
	if !rows.Next() {
		return nil
	}
	return utils.SQLParseRow(rows, u)
}

type Users []*User

func (us *Users) FromSQL(rows *sql.Rows) error {
	for rows.Next() {
		u := &User{}
		if err := utils.SQLParseRow(rows, u); err != nil {
			return err
		}
		*us = append(*us, u)
	}
	return nil
}
