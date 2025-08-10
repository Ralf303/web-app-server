package server

import (
	"example.com/myapp/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
)

func Routes(db *sqlx.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/user/{chatId}", handlers.GetUserHandler(db))
	r.Post("/case/open/{chatId}", handlers.OpenCaseHandler(db))

	return r
}
