package utils

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func SQLParseRow(rows *sql.Rows, into interface{}) error {
	v := reflect.ValueOf(into)
	fields := make([]interface{}, v.Elem().NumField())
	for i := 0; i < v.Elem().NumField(); i++ {
		fields[i] = reflect.New(v.Elem().Field(i).Type()).Interface()
	}
	if err := rows.Scan(fields...); err != nil {
		return err
	}
	for i := 0; i < len(fields); i++ {
		v.Elem().Field(i).Set(reflect.ValueOf(fields[i]).Elem())
	}
	return nil
}

func SQLInsert(ctx context.Context, item interface{}, table string, tx *sql.Tx) error {
	cmd := `INSERT INTO ` + table
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
	_, err := tx.ExecContext(ctx, cmd, fValues...)
	return err
}
