package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/helixauth/helix/cfg"
	"github.com/helixauth/helix/src/shared/utils"
)

func (g *gateway) Txn(ctx context.Context) (Txn, error) {
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

type Txn interface {
	Insert(ctx context.Context, item utils.SQLWritable) error
	Rollback() error
	Commit() error
}

type txn struct {
	tx *sql.Tx
}

func (txn *txn) Insert(ctx context.Context, item utils.SQLWritable) error {
	cmd := `INSERT INTO ` + item.SQLTable()
	fNames := []string{}
	fPlaceholders := []string{}
	fValues := []interface{}{}
	v := reflect.ValueOf(item)
	t := v.Elem().Type()
	for i := 0; i < v.Elem().NumField(); i++ {
		if tag := t.Field(i).Tag.Get("json"); tag != "" {
			fNames = append(fNames, tag)
			fPlaceholders = append(fPlaceholders, fmt.Sprintf("$%v", len(fPlaceholders)+1))
			fValues = append(fValues, v.Elem().Field(i).Interface())
		}
	}
	cmd += fmt.Sprintf(" (%v)", strings.Join(fNames, ", "))
	cmd += fmt.Sprintf(" VALUES (%v)", strings.Join(fPlaceholders, ", "))
	_, err := txn.tx.ExecContext(ctx, cmd, fValues...)
	return err

}

func (txn *txn) Rollback() error {
	return txn.tx.Rollback()
}

func (txn *txn) Commit() error {
	return txn.tx.Commit()
}
