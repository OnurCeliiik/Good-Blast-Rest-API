package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"good-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global database instance
var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	// Set default values if environment variables are missing
	host := getEnv("DB_HOST", "match3-postgres") // Matches docker-compose service name
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password") // Match Docker Compose
	dbname := getEnv("DB_NAME", "match3_db")
	port := getEnv("DB_PORT", "5432") // Default PostgreSQL port

	fmt.Println("Connecting to database with:")
	fmt.Printf("HOST: %s, USER: %s, PASSWORD: %s, DB: %s, PORT: %s\n", host, user, password, dbname, port)

	// Create the DSN (Database Source Name)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable Timezone=UTC",
		host, user, password, dbname, port,
	)

	log.Println("Connecting to database with DSN:", dsn)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}
	if db == nil {
		return nil, errors.New("database connection is nil after initialization")
	}

	log.Println("Connected to database successfully!")

	// Enable uuid-ossp extension (required for generating UUIDs in postgres)
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Printf("Warning: could not enable uuid-ossp extension: %v", err)
	}

	// AutoMigrate will create the table if it does not exist
	err = db.AutoMigrate(&models.User{}, &models.Tournament{}, &models.TournamentParticipant{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return nil, err
	}

	log.Println("Database migration completed!")

	// Assign to global variable
	DB = db
	return db, nil
}

// Helper function to check environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}
