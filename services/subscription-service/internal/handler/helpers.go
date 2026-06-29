package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func getUserID(r *http.Request) uuid.UUID {
	idStr := r.Header.Get("X-User-ID")
	if idStr == "" {
		return uuid.Nil
	}
	id, _ := uuid.Parse(idStr)
	return id
}
