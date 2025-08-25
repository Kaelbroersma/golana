package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func MakeJWT(userid string, tokenSecret string, duration time.Duration) (string, error) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "golana_server",
		Subject:   userid,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	signedToken, err := newToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString string, tokenSecret string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method in JWT")
		}
		return []byte(tokenSecret), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	return claims.Subject, nil
}

func GetBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return "", fmt.Errorf("no bearer token in request")
	}

	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || splitToken[0] != "Bearer" {
		return "", fmt.Errorf("invalid bearer token in request")
	}

	bearerToken := splitToken[1]

	return bearerToken, nil
}
