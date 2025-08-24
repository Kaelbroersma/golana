package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kaelbroersma/golana/internal/database"
	"github.com/Kaelbroersma/golana/internal/server/auth"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ServerConfig struct {
	APIKey       string
	MinVolume    float64
	MinMarketCap float64
	DB           *database.Queries
	Port         string
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Methods for server config

func (cfg *ServerConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	req := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, 500, "Failed to decode request body", err)
		return
	}
	defer r.Body.Close()

	if req.Email == "" {
		RespondWithError(w, 400, "Please provide an email", nil)
		return
	}

	if !auth.IsValidEmail(req.Email) {
		RespondWithError(w, 400, "Email is not valid", nil)
		return
	}

	if req.Password == "" {
		RespondWithError(w, 400, "Password is required", nil)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(w, 500, "Failed to hash password", err)
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:             uuid.New().String(),
		Email:          req.Email,
		Name:           req.Name,
		HashedPassword: string(hashedPassword),
	})
	if err != nil {
		RespondWithError(w, 500, "Failed to create user", err)
		return
	}

	fmt.Printf("Created user: %v\n", req.Name)

	RespondWithJSON(w, http.StatusOK, UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (cfg *ServerConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := LoginUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, 500, "Failed to decode request body", err)
		return
	}
	defer r.Body.Close()

	user, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		RespondWithError(w, 500, "Failed to get user-check username and try again.", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		RespondWithError(w, 500, "Invalid password", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// Handlers that don't require config

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, "Hello, World!")
}
