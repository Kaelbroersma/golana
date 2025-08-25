package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to marshal response", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Println("Failed to write response", err)
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string, err error) {
	errorToRespond := errorResponse{
		Error: message,
	}
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Println("Responding with 5xx error:", err)
	}

	RespondWithJSON(w, code, errorToRespond)
}
