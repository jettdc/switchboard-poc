package u

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func InitializeEnv(envFile string) error {
	err := godotenv.Load(envFile)
	if err != nil {
		return err
	}
	return nil
}

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
		msg := fmt.Sprintf("env lookup for value \"%s\" failed. Using default value: \"%s\"", name, defaultValue)
		if Logger != nil {
			Logger.Warn(msg)
		}
		return defaultValue
	}
	return res
}
