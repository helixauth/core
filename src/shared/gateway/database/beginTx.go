package database

import (
	"context"
	"database/sql"
	"fmt"
)

func (g *gateway) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	tenantID, _ := ctx.Value("TENANT_ID").(string)
	cmd := fmt.Sprintf("SET app.tenant_id = '%v';", tenantID)
	if _, err := tx.ExecContext(ctx, cmd); err != nil {
		return nil, err
	}

	return tx, nil
}
