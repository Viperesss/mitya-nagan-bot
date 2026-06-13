package e

import "fmt"

func Wrap(message string, err error) error {
	return fmt.Errorf("%s %w", message, err)
}
