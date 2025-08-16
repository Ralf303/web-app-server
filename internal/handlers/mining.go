package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/myapp/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type SlotResponse struct {
	Id   int            `json:"id"`
	Card *database.Card `json:"card"`
}

type InstallGpuRequest struct {
	UserId  string `json:"userId"`
	StandId int    `json:"standId"`
	CardId  int    `json:"cardId"`
}

type BuyGpuSlotRequest struct {
	Status string             `json:"status"`
	Slot   database.CardStand `json:"slot"`
}

func GetSlotsHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr := chi.URLParam(r, "userId")

		user, err := database.GetUser(db, userIdStr)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		stands, err := database.GetUserCardStands(db, user.Id)
		if err != nil {
			http.Error(w, "Failed to get slots", http.StatusInternalServerError)
			return
		}

		slots := make([]SlotResponse, 0, len(stands))
		for _, stand := range stands {
			slots = append(slots, SlotResponse{
				Id:   stand.Id,
				Card: stand.Card,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(slots); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func GetGpuHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr := chi.URLParam(r, "userId")

		user, err := database.GetUser(db, userIdStr)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		cards, err := database.GetUserCards(db, user.Id)
		if err != nil {
			http.Error(w, "Failed to get GPUs", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cards)
	}
}

func GetGpuByIdHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gpuIdStr := chi.URLParam(r, "gpuId")
		gpuId, err := strconv.Atoi(gpuIdStr)
		if err != nil {
			http.Error(w, "Invalid gpuId", http.StatusBadRequest)
			return
		}

		card, err := database.GetCardById(db, gpuId)
		if err != nil {
			http.Error(w, "GPU not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(card)
	}
}

func InstallGpuHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req InstallGpuRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := database.GetUser(db, req.UserId)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		card, err := database.GetCardById(db, req.CardId)
		if err != nil {
			http.Error(w, "Card not found", http.StatusNotFound)
			return
		}
		if card.UserId != user.Id {
			http.Error(w, "Card does not belong to user", http.StatusForbidden)
			return
		}

		stand, err := database.GetCardStandById(db, req.StandId)
		if err != nil {
			http.Error(w, "Stand not found", http.StatusNotFound)
			return
		}
		if stand.UserId != user.Id {
			http.Error(w, "Stand does not belong to user", http.StatusForbidden)
			return
		}

		installedElsewhere, err := database.IsCardInstalledElsewhere(db, req.CardId)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if installedElsewhere {
			http.Error(w, "Card is already installed in another stand", http.StatusConflict)
			return
		}

		if stand.CardId != nil {
			http.Error(w, "Stand already has a card", http.StatusConflict)
			return
		}

		newCoins := user.Coin + card.Balance
		if err := database.ResetCardBalance(db, card.Id); err != nil {
			http.Error(w, "Failed to update card balance", http.StatusInternalServerError)
			return
		}
		if err := database.UpdateUserCoins(db, user.Id, newCoins); err != nil {
			http.Error(w, "Failed to update user coins", http.StatusInternalServerError)
			return
		}
		if err := database.InsertCardIntoStand(db, stand.Id, card.Id); err != nil {
			http.Error(w, "Failed to install card into stand", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}

func BuySlotHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr := chi.URLParam(r, "userId")

		user, err := database.GetUser(db, userIdStr)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		stands, err := database.GetUserCardStands(db, user.Id)
		if err != nil {
			http.Error(w, "Failed to get slots", http.StatusInternalServerError)
			return
		}
		if len(stands) >= 9 {
			http.Error(w, "maxSlots", http.StatusBadRequest)
			return
		}

		const slotPrice = 2500000
		if user.Balance < uint64(slotPrice) {
			http.Error(w, "dontMoney", http.StatusBadRequest)
			return
		}

		newBalance := user.Balance - uint64(slotPrice)
		err = database.UpdateUserBalance(db, user.Id, newBalance)
		if err != nil {
			http.Error(w, "Failed to update balance", http.StatusInternalServerError)
			return
		}

		stand, err := database.CreateCardStand(db, user.Id)
		if err != nil {
			http.Error(w, "Failed to create slot", http.StatusInternalServerError)
			return
		}

		response := BuyGpuSlotRequest{
			Status: "success",
			Slot:   stand,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
