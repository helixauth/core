package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/shared/utils"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type Gateway interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Query(ctx context.Context, into utils.SQLReadable, qry string, args ...interface{}) error
}

type gateway struct {
	db *sql.DB
}

func New(ctx context.Context) (Gateway, error) {
	connInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		ctx.Value(cfg.PostgresHost).(string),
		ctx.Value(cfg.PostgresPort).(string),
		ctx.Value(cfg.PostgresUsername).(string),
		ctx.Value(cfg.PostgresPassword).(string),
		ctx.Value(cfg.PostgresDBName).(string),
		ctx.Value(cfg.PostgresSSLMode).(string),
	)
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to database")
	}
	return &gateway{
		db: db,
	}, err
}
