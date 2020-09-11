package mapper

// BoolPtr maps a bool to a bool pointer
func BoolPtr(b bool) *bool {
	return &b
}

// StringPtr maps a string to a string pointer
func StringPtr(s string) *string {
	return &s
}
