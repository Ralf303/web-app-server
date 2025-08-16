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
	r.Get("/mining/withdrowBitcoin/{cardId}/{userId}", handlers.WithdrawBitcoinHandler(db))
	r.Get("/mining/getSlots/{userId}", handlers.GetSlotsHandler(db))
	r.Get("/mining/getGpu/{userId}", handlers.GetGpuHandler(db))
	r.Get("/mining/getGpuById/{gpuId}", handlers.GetGpuByIdHandler(db))
	r.Get("/mining/pullGpu/{gpuId}/{userId}", handlers.PullGpuHandler(db))
	r.Post("/case/open/{chatId}", handlers.OpenCaseHandler(db))
	r.Post("/mining/installGpu", handlers.InstallGpuHandler(db))
	r.Post("/mining/buySlot/{userId}", handlers.BuySlotHandler(db))
	r.Post("/mining/freezeGpu", handlers.FreezeGpuHandler(db))

	return r
}
