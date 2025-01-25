package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"eventix/entity"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBInit() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)

	var db *gorm.DB
	for i := 0; i < 10; i++ { // Retry 10 kali
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database. Retrying in 5 seconds... (%d/10)\n", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Automigrate tables
	err = db.AutoMigrate(
		&entity.Event{},
		&entity.Ticket{},
		&entity.User{},
		&entity.TokenBlacklist{},
		// Add other entities here
	)
	if err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	return db
}
