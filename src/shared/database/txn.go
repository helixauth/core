package database

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
)

type Writable interface {
	DatabaseTable() string
}

type Txn interface {
	Insert(ctx context.Context, item Writable) error
	Rollback() error
	Commit() error
}

type txn struct {
	tx *sql.Tx
}

func (txn *txn) Insert(ctx context.Context, item Writable) error {
	cmd := `INSERT INTO ` + item.DatabaseTable()
	fNames := []string{}
	fPlaceholders := []string{}
	fValues := []interface{}{}
	v := reflect.ValueOf(item)
	t := v.Elem().Type()
	for i := 0; i < v.Elem().NumField(); i++ {
		if tag := t.Field(i).Tag.Get("json"); tag != "" {
			fNames = append(fNames, tag)
			fPlaceholders = append(fPlaceholders, fmt.Sprintf("$%v", len(fPlaceholders)+1))
			switch v.Elem().Field(i).Kind() {
			case reflect.Slice:
				fValues = append(fValues, pq.Array(v.Elem().Field(i).Interface()))
			default:
				fValues = append(fValues, v.Elem().Field(i).Interface())
			}
		}
	}
	cmd += fmt.Sprintf(" (%v)", strings.Join(fNames, ", "))
	cmd += fmt.Sprintf(" VALUES (%v)", strings.Join(fPlaceholders, ", "))
	_, err := txn.tx.ExecContext(ctx, cmd, fValues...)
	if err != nil {
		txn.Rollback()
	}
	return err
}

func (txn *txn) Rollback() error {
	return txn.tx.Rollback()
}

func (txn *txn) Commit() error {
	return txn.tx.Commit()
}
