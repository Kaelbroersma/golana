package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

func StartServer(cfg *Config) {
	if cfg.DB == nil {
		fmt.Println("DATABASE_URL is not set or database connection failed.")
		fmt.Println("Running without CRUD operations")
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", cfg.Port),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
	}

	mux.Handle("GET /api", http.HandlerFunc(handleGetRoot))
	mux.Handle("GET /api/health", http.HandlerFunc(handleReadiness))

	mux.Handle("POST /api/users", http.HandlerFunc(cfg.handleCreateUser))
	mux.Handle("POST /api/login", http.HandlerFunc(cfg.handleLogin))
	mux.Handle("POST /api/trades", http.HandlerFunc(cfg.handleCreateTrade))

	mux.Handle("UPDATE /api/trades", http.HandlerFunc(cfg.handleCloseTrade))

	fmt.Printf("Server is running on port %v\n", cfg.Port)

	log.Fatal(srv.ListenAndServe())

}
