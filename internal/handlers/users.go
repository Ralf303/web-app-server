package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/myapp/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func GetUserHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId := chi.URLParam(r, "chatId")
		user, err := database.GetUser(db, chatId)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}
