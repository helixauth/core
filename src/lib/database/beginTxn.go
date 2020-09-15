package database

import (
	"context"
	"fmt"

	"github.com/helixauth/core/cfg"
)

// BeginTxn begins a new database transaction
func (g *gateway) BeginTxn(ctx context.Context) (Txn, error) {
	tx, err := g.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	tenantID := ctx.Value(cfg.TenantID).(string)
	cmd := fmt.Sprintf("SET app.tenant_id = '%v';", tenantID)
	if _, err := tx.ExecContext(ctx, cmd); err != nil {
		return nil, err
	}
	return &txn{
		tx: tx,
	}, nil
}
