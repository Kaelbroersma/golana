package server

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/tursodatabase/go-libsql"
)

func StartServer(cfg *ServerConfig) {
	if cfg.DB == nil {
		fmt.Println("DATABASE_URL is not set or database connection failed.")
		fmt.Println("Running without CRUD operations")
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: mux,
	}

	mux.Handle("GET /api", http.HandlerFunc(handleGetRoot))
	mux.Handle("GET /api/health", http.HandlerFunc(handleReadiness))

	mux.Handle("POST /api/users", http.HandlerFunc(cfg.handleCreateUser))
	mux.Handle("POST /api/login", http.HandlerFunc(cfg.handleLogin))
	mux.Handle("POST /api/trades", http.HandlerFunc(cfg.handleCreateTrade))

	fmt.Printf("Server is running on port %v\n", cfg.Port)

	log.Fatal(srv.ListenAndServe())

}
