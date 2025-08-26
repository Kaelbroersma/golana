package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Kaelbroersma/golana/internal/database"
	"github.com/Kaelbroersma/golana/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	cfg := server.Config{}

	// Load environment

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, is it present?")
	}

	// Open a connection to the database

	db, err := sql.Open("libsql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	cfg.DB = dbQueries

	// Get port

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	cfg.Port = port

	// Get token secret - WE RECOMMEND USING AN ENCRYPTED TOKEN SECRET. AVOID USING SOMETHING DETERMINISTIC.

	tokenSecret := os.Getenv("TOKEN_SECRET")
	if tokenSecret == "" {
		log.Fatal("TOKEN_SECRET is not set")
	}
	cfg.TokenSecret = tokenSecret

	// Get HELIUS API KEY

	heliusAPIKey := os.Getenv("HELIUS_API_KEY")
	if heliusAPIKey == "" {
		log.Fatal("HELIUS_API_KEY is not set")
	}
	cfg.HeliusAPIKey = heliusAPIKey

	// Start Server

	go server.StartServer(&cfg)

	// Keep server running

	select {}

}
