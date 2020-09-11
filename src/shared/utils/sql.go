package utils

import (
	"database/sql"
	"reflect"
)

type SQLWritable interface {
	SQLTable() string
}

type SQLReadable interface {
	FromSQL(*sql.Rows) error
}

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
