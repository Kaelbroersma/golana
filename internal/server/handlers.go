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

type Config struct {
	APIKey       string
	MinVolume    float64
	MinMarketCap float64
	DB           *database.Queries
	Port         string
	TokenSecret  string
	HeliusAPIKey string
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
	Token     string    `json:"token"`
	Name      string    `json:"name"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Methods for server config

func (cfg *Config) handleCreateUser(w http.ResponseWriter, r *http.Request) {
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

	userToken, err := auth.MakeJWT(user.ID, cfg.TokenSecret, 24*time.Hour)
	if err != nil {
		RespondWithError(w, 500, "Failed to create JWT", err)
		return
	}

	fmt.Printf("Created user: %v\n", req.Name)

	RespondWithJSON(w, http.StatusOK, UserResponse{
		Token:     userToken,
		Name:      user.Name,
		Balance:   user.BuyingPower,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (cfg *Config) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := LoginUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, 500, "Failed to decode request body", err)
		return
	}
	defer r.Body.Close()

	user, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		RespondWithError(w, 500, "Could not find user. Please try again.", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		RespondWithError(w, 500, "Invalid password", err)
		return
	}

	userToken, err := auth.MakeJWT(user.ID, cfg.TokenSecret, 30*time.Minute)
	if err != nil {
		RespondWithError(w, 500, "Failed to create JWT", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, UserResponse{
		Token:     userToken,
		Name:      user.Name,
		Balance:   user.BuyingPower,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// Handlers that don't require config

func handleGetRoot(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, "Hello, World!")
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, "OK")
}
