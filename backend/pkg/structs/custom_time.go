// Package structs defines custom data structures used throughout the application.
package structs

import (
	"fmt"
	"strings"
	"time"
)

// CustomTime is a wrapper around time.Time to support custom JSON unmarshaling.
type CustomTime struct {
	time.Time
}

// UnmarshalJSON unmarshals JSON data into CustomTime, supporting RFC3339 and YYYY-MM-DD formats.
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}

	// Try RFC3339
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		ct.Time = t
		return nil
	}

	// Try YYYY-MM-DD
	t, err = time.Parse("2006-01-02", s)
	if err == nil {
		ct.Time = t
		return nil
	}

	return fmt.Errorf("invalid time format: %s", s)
}

// MarshalJSON marshals CustomTime into JSON data in RFC3339 format.
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	if ct.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ct.Format(time.RFC3339))), nil
}
