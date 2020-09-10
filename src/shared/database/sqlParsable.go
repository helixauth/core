package database

import (
	"database/sql"
)

type SQLParsable interface {
	FromSQL(*sql.Rows) error
}
