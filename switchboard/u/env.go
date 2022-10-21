package u

import (
	"fmt"
	"os"
)

func ValidateRequiredEnv(variables []string) error {
	for _, variable := range variables {
		if _, exists := os.LookupEnv(variable); !exists {
			return fmt.Errorf("missing required environment variable %s", variable)
		}
	}
	return nil
}

func GetEnvWithDefault(name, defaultValue string) string {
	res, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}
	return res
}
