package database

import (
	"context"
	"fmt"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/shared/utils"
)

func (g *gateway) Query(ctx context.Context, into utils.SQLReadable, qry string, args ...interface{}) error {
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

	rows, err := conn.QueryContext(ctx, qry, args...)
	if err != nil {
		return err
	}

	defer rows.Close()
	return into.FromSQL(rows)
}
