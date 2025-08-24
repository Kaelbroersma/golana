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
	cfg := server.ServerConfig{}

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	cfg.Port = port

	// Start Server

	go server.StartServer(&cfg)

	// Keep server running

	select {}

}
