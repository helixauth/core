package database

import (
	"context"
	"database/sql"
)

func (g *gateway) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return g.db.BeginTx(ctx, nil)
}
