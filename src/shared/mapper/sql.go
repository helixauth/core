package mapper

import (
	"database/sql"
	"encoding/json"
)

// StringPtrToSQLNullString maps a *string to a sql.NullString
func StringPtrToSQLNullString(str *string) sql.NullString {
	s := ""
	if str != nil {
		s = *str
	}
	return sql.NullString{
		String: s,
		Valid:  str != nil,
	}
}

// MapToSQLJSON maps a map[string]interface{} to an SQL JSON value
func MapToSQLJSON(m interface{}) string {
	bytes, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(bytes)
}

// SQLNullStringToStringPtr maps a sql.NullString to a *string
func SQLNullStringToStringPtr(str sql.NullString) *string {
	if !str.Valid {
		return nil
	}
	return &str.String
}
