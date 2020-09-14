package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/helixauth/helix/src/lib/secrets"

	"github.com/pkg/errors"
)

// Gateway provides an interface to the PostgreSQL database
type Gateway interface {
	BeginTxn(ctx context.Context) (Txn, error)
	QueryItem(ctx context.Context, item interface{}, qry string, args ...interface{}) error
	QueryList(ctx context.Context, list interface{}, qry string, args ...interface{}) error
}

type gateway struct {
	db *sql.DB
}

// New creates a new database gateway
func New(ctx context.Context, secretsManager secrets.Manager) (Gateway, error) {
	var err error
	connArgs := map[string]interface{}{
		"postgres.host":     "",
		"postgres.port":     "",
		"postgres.username": "",
		"postgres.password": "",
		"postgres.db_name":  "",
		"postgres.ssl_mode": "",
	}
	for k := range connArgs {
		if connArgs[k], err = secretsManager.Get(k); err != nil {
			return nil, err
		}
	}
	connInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		connArgs["postgres.host"],
		connArgs["postgres.port"],
		connArgs["postgres.username"],
		connArgs["postgres.password"],
		connArgs["postgres.db_name"],
		connArgs["postgres.ssl_mode"],
	)
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to database")
	}
	return &gateway{
		db: db,
	}, err
}
