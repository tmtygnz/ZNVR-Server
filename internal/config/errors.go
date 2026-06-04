package config

import "fmt"

func MissingFieldError(missing string) error {
	return fmt.Errorf("Missing or wrong type: %s", missing)
}
