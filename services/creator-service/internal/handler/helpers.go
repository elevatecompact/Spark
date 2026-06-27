package handler

import (
	"encoding/json"
	"net/http"

	"github.com/elevatecompact/spark/services/creator-service/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeDomainError(w http.ResponseWriter, err error) {
	status := domain.MapErrorToStatus(err)
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
