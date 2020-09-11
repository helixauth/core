package database

import (
	"database/sql"
	"fmt"
	"reflect"
)

func parseRows(rows *sql.Rows, list interface{}) error {

	// Validate 'list' is a pointer to a slice
	vListPtr := reflect.ValueOf(list)
	if vListPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("'list' must be a pointer to a slice")
	}
	vList := vListPtr.Elem()
	if vList.Kind() != reflect.Slice {
		return fmt.Errorf("'list' must be a pointer to a slice")
	}

	// Get the item type of the slice
	tItem := vList.Type().Elem()

	// For each SQL row...
	for rows.Next() {

		// Generate a pointer to a new item
		vItemPtr := reflect.New(tItem)

		// Scan the SQL row into the item
		if err := scan(rows, vItemPtr.Elem()); err != nil {
			return err
		}

		// Append the item to the list
		vList.Set(reflect.Append(vList, vItemPtr.Elem()))
	}

	return nil
}

func parseRow(row *sql.Row, item interface{}) error {

	// Validate 'item' is a pointer to a struct
	vItemPtr := reflect.ValueOf(item)
	if vItemPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("'item' must be a pointer to a struct")
	}
	vItem := vItemPtr.Elem()
	if vItem.Kind() != reflect.Struct {
		return fmt.Errorf("'item' must be a pointer to a struct")
	}

	// Scan the SQL row into the item
	return scan(row, vItem)
}

type scannable interface {
	Scan(dest ...interface{}) error
}

func scan(s scannable, vItem reflect.Value) error {

	// Create slice of pointers
	// Each pointer's value type matches a field type of the struct being written to
	ptrs := make([]interface{}, vItem.NumField())
	for i := 0; i < vItem.NumField(); i++ {
		ptrs[i] = reflect.New(vItem.Field(i).Type()).Interface()
	}

	// Scan row in the pointer values
	if err := s.Scan(ptrs...); err != nil {
		return err
	}

	// Assign the pointer values to struct fields
	for i := 0; i < len(ptrs); i++ {
		vItem.Field(i).Set(reflect.ValueOf(ptrs[i]).Elem())
	}

	return nil
}
