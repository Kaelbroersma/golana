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

	mux.Handle("GET /", http.HandlerFunc(handleGetRoot))
	mux.Handle("POST /users", http.HandlerFunc(cfg.handleCreateUser))
	mux.Handle("POST /login", http.HandlerFunc(cfg.handleLogin))

	fmt.Printf("Server is running on port %v\n", cfg.Port)

	log.Fatal(srv.ListenAndServe())

}
