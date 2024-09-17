package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/myapp/internal/database"
	user "example.com/myapp/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
)

func Routes(db *sqlx.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://pablohouse.su"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/user/get/{userId}", func(w http.ResponseWriter, r *http.Request) {
		GetUser(db, w, r)
	})
	router.Put("/user/updateBalance/{userId}", func(w http.ResponseWriter, r *http.Request) {
		UpdateUserBalance(db, w, r)
	})
	router.Put("/user/updateGems/{userId}", func(w http.ResponseWriter, r *http.Request) {
		UpdateUserGems(db, w, r)
	})
	router.Put("/user/updateKeys/{userId}", func(w http.ResponseWriter, r *http.Request) {
		UpdateUserKeys(db, w, r)
	})
	return router
}

func GetUser(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	user, err := database.GetOrCreateUser(db, userId)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error serializing user data", 500)
		return
	}

	w.Write(userJSON)
}

func UpdateUserBalance(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	var request struct {
		Balance uint64 `json:"balance"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := user.UpdatedBalance(db, userId, request.Balance)
	if err != nil {
		http.Error(w, "Error updating balance", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User %s balance updated to %d\n", userId, request.Balance)
}

func UpdateUserGems(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	var request struct {
		Gems uint64 `json:"gems"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := user.UpdatedGems(db, userId, request.Gems)
	if err != nil {
		http.Error(w, "Error updating gems", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User %s gems updated to %d\n", userId, request.Gems)
}

func UpdateUserKeys(db *sqlx.DB, w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	var request struct {
		Keys uint64 `json:"keys"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := user.UpdatedKeys(db, userId, request.Keys)
	if err != nil {
		http.Error(w, "Error updating keys", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User %s keys updated to %d\n", userId, request.Keys)
}
