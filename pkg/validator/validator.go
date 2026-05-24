package validator

import "strings"

// Validate checks that required string fields are non-empty.
func Validate(fields map[string]string) map[string]string {
	errors := make(map[string]string)
	for field, value := range fields {
		if strings.TrimSpace(value) == "" {
			errors[field] = field + " is required"
		}
	}
	if len(errors) == 0 {
		return nil
	}
	return errors
}
