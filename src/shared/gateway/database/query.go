package database

import (
	"context"
)

func (g *gateway) Query(ctx context.Context, into SQLParsable, qry string, args ...interface{}) error {
	rows, err := g.db.Query(qry, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	return into.FromSQL(rows)
}
