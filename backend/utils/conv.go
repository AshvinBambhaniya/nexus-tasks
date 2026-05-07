// Package utils provides various utility functions.
package utils

import "fmt"

// GetString converts any type to string using fmt.Sprint
func GetString(text any) string {
	return fmt.Sprint(text)
}
