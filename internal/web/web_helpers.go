package web

import "time"

func optional(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func defaultEnrollmentDate(value string) string {
	if value != "" {
		return value
	}
	return time.Now().Format(time.DateOnly)
}
