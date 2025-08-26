package auth

import (
	"fmt"
	"net/http"
	"net/mail"

	"github.com/Kaelbroersma/golana/internal/database"
	"github.com/Kaelbroersma/golana/internal/server"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func AuthenticateWithBearer(r *http.Request, cfg *server.ServerConfig) (database.User, error) {
	bearerToken, err := GetBearerToken(r)
	if err != nil {
		return database.User{}, err
	}

	userID, err := ValidateJWT(bearerToken, cfg.TokenSecret)
	if err != nil {
		return database.User{}, err
	}

	user, err := cfg.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		return database.User{}, fmt.Errorf("could not get user")
	}

	return user, nil
}
