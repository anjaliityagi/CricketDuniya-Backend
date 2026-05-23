package main

import (
	databasepackage "CricketDuniya-Backend/internal/database"
	"CricketDuniya-Backend/internal/server"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	err = databasepackage.ConnectAndMigrate(
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPassword,
		databasepackage.SSLMode(sslMode),
	)

	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("Database connected successfully")

	srv := server.SetupRoutes()

	fmt.Println("Server running at http://172.16.2.37:8080")

	if err := srv.Router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
