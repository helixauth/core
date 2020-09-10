package database

import (
	"context"
	"fmt"
)

func (g *gateway) Query(ctx context.Context, into SQLParsable, qry string, args ...interface{}) error {
	conn, err := g.db.Conn(ctx)
	if err != nil {
		conn.Close()
		return err
	}

	tenantID, _ := ctx.Value("TENANT_ID").(string)
	cmd := fmt.Sprintf("SET app.tenant_id = '%v';", tenantID)
	if _, err = conn.ExecContext(ctx, cmd); err != nil {
		conn.Close()
		return err
	}

	rows, err := conn.QueryContext(ctx, qry, args...)
	if err != nil {
		return err
	}

	defer rows.Close()
	return into.FromSQL(rows)
}
