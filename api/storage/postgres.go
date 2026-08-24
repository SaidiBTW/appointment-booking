package storage

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/SaidiBTW/appointment_booking_system_go/config"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.PostgresConfig) (*sql.DB, error) {
	connStr := "host=" + cfg.Host + " port=" + strconv.Itoa(cfg.Port) + " user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.DBName + " sslmode=" + cfg.SSLMode
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		db.Close()
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "development" {
		log.Println("Seeding the database with initial data...")
		SeedDatabase(db)
		log.Println("Database seeding completed.")
	}

	log.Println("Connected to PostgreSQL database successfully")
	return db, nil
}
