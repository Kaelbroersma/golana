package auth

import (
	"fmt"
	"net/http"
	"net/mail"

	"github.com/Kaelbroersma/golana/internal/database"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func AuthenticateWithBearer(r *http.Request, tokenSecret string, db *database.Queries) (database.User, error) {
	bearerToken, err := GetBearerToken(r)
	if err != nil {
		return database.User{}, err
	}

	userID, err := ValidateJWT(bearerToken, tokenSecret)
	if err != nil {
		return database.User{}, err
	}

	user, err := db.GetUserByID(r.Context(), userID)
	if err != nil {
		return database.User{}, fmt.Errorf("could not get user")
	}

	return user, nil
}
