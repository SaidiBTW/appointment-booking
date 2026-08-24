package config

import (
	"os"
	"strconv"
)

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

func DefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "appointment_booking",
		SSLMode:  "disable",
	}
}

func Load() (*PostgresConfig, error) {
	// Load configuration from environment variables or a config file
	// For simplicity, we will return the default configuration here
	defaultConfig := DefaultPostgresConfig()
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	sslmode := os.Getenv("DATABASE_SSLMODE")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		return defaultConfig, nil
	}

	if sslmode == "" {
		sslmode = defaultConfig.SSLMode
	}

	return &PostgresConfig{
		Host:     host,
		Port:     func() int { p, _ := strconv.Atoi(port); return p }(),
		User:     user,
		Password: password,
		DBName:   dbname,
		SSLMode:  sslmode,
	}, nil
}

func ReadSeederSQL() (string, error) {
	// Read the contents of the ./seeder.sql file
	filePath := "./seeder.sql"
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
