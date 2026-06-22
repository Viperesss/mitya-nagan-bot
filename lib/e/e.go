// Package e provides helpers for error handling.
package e

import "fmt"

// Wrap wraps an error with an additional message.
func Wrap(message string, err error) error {
	return fmt.Errorf("%s %w", message, err)
}
