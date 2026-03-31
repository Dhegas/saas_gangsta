package main

import (
	"log"

	"github.com/dhegas/saas_gangsta/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	database.Connect()

	// router...
}
