package database

import (
	"context"
	"fmt"

	"github.com/helixauth/helix/cfg"
)

func (g *gateway) QueryItem(ctx context.Context, item interface{}, qry string, args ...interface{}) error {
	conn, err := g.db.Conn(ctx)
	if err != nil {
		return err
	}

	tenantID := ctx.Value(cfg.TenantID).(string)
	cmd := fmt.Sprintf("SET app.tenant_id = '%v';", tenantID)
	if _, err = conn.ExecContext(ctx, cmd); err != nil {
		conn.Close()
		return err
	}

	row := conn.QueryRowContext(ctx, qry, args...)
	return parseRow(row, item)
}
