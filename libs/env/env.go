package env

import (
	"fmt"
	"os"
	"strconv"
)

var readers []reader

type reader func() error

func RequiredStringVar(v *string, name string) {
	readers = append(readers, func() error {
		value := os.Getenv(name)
		if value == "" {
			return fmt.Errorf("environment variable %s is not set", name)
		}

		*v = value

		return nil
	})
}

func StringVar(v *string, name string, otherwise string) {
	readers = append(readers, func() error {
		if value := os.Getenv(name); value == "" {
			*v = otherwise
		} else {
			*v = value
		}

		return nil
	})
}

func BoolVar(v *bool, name string, otherwise bool) {
	readers = append(readers, func() error {
		stringValue := os.Getenv(name)
		if stringValue == "" {
			*v = otherwise
			return nil
		}

		boolValue, err := strconv.ParseBool(stringValue)
		if err != nil {
			return fmt.Errorf("environment variable %s is not a boolean", name)
		}

		*v = boolValue

		return nil
	})
}

func Parse() error {
	defer func() {
		readers = nil
	}()

	for _, reader := range readers {
		if err := reader(); err != nil {
			return err
		}
	}

	return nil
}
