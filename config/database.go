package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
func ConnectDB() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Check if environment variables are set
	if dbUser == "" || dbName == "" || dbHost == "" || dbPort == "" {
		log.Fatal("Database environment variables are not set properly")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}

	log.Println("Database connected")
	DB = db
}