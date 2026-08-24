package config

import (
	"log"
	"os"
	"strconv"
)

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

func DefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "appointment_booking",
	}
}

func Load() (*PostgresConfig, error) {
	// Load configuration from environment variables or a config file
	// For simplicity, we will return the default configuration here
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Println("Environment variables for database configuration are not set. Using default configuration.")
		return DefaultPostgresConfig(), nil
	}

	return &PostgresConfig{
		Host:     host,
		Port:     func() int { p, _ := strconv.Atoi(port); return p }(),
		User:     user,
		Password: password,
		DBName:   dbname,
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
