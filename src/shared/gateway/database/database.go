package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/helixauth/helix/cfg"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Gateway interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Query(ctx context.Context, into SQLParsable, qry string, args ...interface{}) error
}

type gateway struct {
	db *sql.DB
}

func New(cfg config.Config) (Gateway, error) {
	connInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.DBName, cfg.Postgres.SSLMode)
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to database")
	}
	return &gateway{
		db: db,
	}, err
}
