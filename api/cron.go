package handler

import (
	"encoding/json"
	"net/http"
	"os"
)

func Cron(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	cronSecret := os.Getenv("CRON_SECRET")

	if authHeader != "Bearer "+cronSecret {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"data": "Hello, Cron!",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
