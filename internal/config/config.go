package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	Environment string
}

func Load() (*Config, error) {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("Error loading .env file: %s\n", err)
	}

	cfg := &Config{
		Port:        getEnvVar("PORT", "8080"),
		DatabaseURL: getEnvVar("DATABASE_URL"),
		Environment: getEnvVar("ENVIRONMENT", "development"),
	}

	if cfg.DatabaseURL == "" {
		// Build from individual components if DATABASE_URL not set
		host := getEnvVar("DB_HOST", "localhost")
		port := getEnvVar("DB_PORT", "5432")
		user := getEnvVar("DB_USER", "postgres")
		password := getEnvVar("DB_PASSWORD", "")
		dbname := getEnvVar("DB_NAME", "discord_resources")
		sslmode := getEnvVar("DB_SSLMODE", "require")
		if sslmode == "disable" {
			log.Println("Warning: Database SSL mode is disabled")
		}

		cfg.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbname, sslmode,
		)
	}

	return cfg, nil
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

/*
* Gets an environment variable, loading from .env file if needed.
* If the variable is not set, it sets it to the provided default value.
 */
func getEnvVar(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	if len(defaultValue) > 0 {
		log.Printf("Using default value for environment variable: %s", key)
		return defaultValue[0]
	}

	log.Printf("Environment variable not set: %s", key)
	return ""
}
