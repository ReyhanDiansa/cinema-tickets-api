package models

import (
	"cinema-tickets/config"
	"log"
)


func AutoMigrate(){
	err := config.DB.AutoMigrate(
		&Film{},
		&Schedule{},
		&Seat{},
		&Transaction{},
		&User{},
		&TransactionSeat{},
		&Cinema{},
	)

	if err != nil {
		log.Fatalf("Error during migration: %v", err)
	}

	log.Println("Auto migration completed.")
}