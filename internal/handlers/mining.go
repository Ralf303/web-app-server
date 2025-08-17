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

type FreezeGpuRequest struct {
	UserId string `json:"userId"`
	CardId int    `json:"cardId"`
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

		if cards == nil {
			cards = []database.Card{}
		}

		type CardWithIncome struct {
			Card   database.Card `json:"card"`
			Income int           `json:"income"`
		}

		var result []CardWithIncome
		for _, c := range cards {
			result = append(result, CardWithIncome{
				Card:   c,
				Income: c.Lvl + 10,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
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
			json.NewEncoder(w).Encode(map[string]string{"status": "dontHaveGpu"})
			return
		}

		stand, err := database.GetCardStandById(db, req.StandId)
		if err != nil {
			http.Error(w, "Stand not found", http.StatusNotFound)
			return
		}

		if stand.UserId != user.Id {
			json.NewEncoder(w).Encode(map[string]string{"status": "dontHaveStand"})
			return
		}

		installedElsewhere, err := database.IsCardInstalledElsewhere(db, req.CardId)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if installedElsewhere {
			json.NewEncoder(w).Encode(map[string]string{"status": "alreadyInstalledElsewhere"})
			return
		}

		if stand.CardId != nil {
			json.NewEncoder(w).Encode(map[string]string{"status": "alreadyInstalled"})
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
			json.NewEncoder(w).Encode(map[string]string{"status": "maxSlots"})
			return
		}

		const slotPrice = 2500000
		if user.Balance < uint64(slotPrice) {
			json.NewEncoder(w).Encode(map[string]string{"status": "noBalance"})
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

func FreezeGpuHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req FreezeGpuRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		user, err := database.GetUser(db, req.UserId)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Freeze <= 0 {
			http.Error(w, "dontFreeze", http.StatusBadRequest)
			return
		}

		card, err := database.GetCardById(db, req.CardId)
		if err != nil {
			http.Error(w, "GPU not found", http.StatusNotFound)
			return
		}

		if card.UserId != user.Id {
			json.NewEncoder(w).Encode(map[string]string{"status": "dontHaveGpu"})
			return
		}

		if card.Fuel >= 100 {
			json.NewEncoder(w).Encode(map[string]string{"status": "alreadyFull"})
			return
		}

		newFuel := min(card.Fuel+50, 100)

		if err := database.UpdateCardFuel(db, card.Id, newFuel); err != nil {
			http.Error(w, "Failed to update GPU fuel", http.StatusInternalServerError)
			return
		}

		if err := database.DecrementUserFreeze(db, user.Id); err != nil {
			http.Error(w, "Failed to update user freeze", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"status": "success",
			"card": map[string]interface{}{
				"id":     card.Id,
				"fuel":   newFuel,
				"userId": card.UserId,
			},
			"user": map[string]interface{}{
				"id":     user.Id,
				"freeze": user.Freeze - 1,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func WithdrawBitcoinHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cardIdStr := chi.URLParam(r, "cardId")
		userIdStr := chi.URLParam(r, "userId")

		cardId, err := strconv.Atoi(cardIdStr)
		if err != nil {
			http.Error(w, "Invalid cardId", http.StatusBadRequest)
			return
		}

		user, err := database.GetUser(db, userIdStr)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		card, err := database.GetCardById(db, cardId)
		if err != nil {
			http.Error(w, "Card not found", http.StatusNotFound)
			return
		}

		if card.UserId != user.Id {
			json.NewEncoder(w).Encode(map[string]string{"status": "dontHave"})
			return
		}

		cardBalance := card.Balance
		if cardBalance <= 0 {
			json.NewEncoder(w).Encode(map[string]string{"status": "noBalance"})
			return
		}

		if err := database.ResetCardBalance(db, card.Id); err != nil {
			http.Error(w, "Failed to reset card balance", http.StatusInternalServerError)
			return
		}

		newUserCoins := user.Coin + cardBalance
		if err := database.UpdateUserCoins(db, user.Id, newUserCoins); err != nil {
			http.Error(w, "Failed to update user balance", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"status": "success",
			"user": map[string]interface{}{
				"id":   user.Id,
				"coin": newUserCoins,
			},
			"card": map[string]interface{}{
				"id":      card.Id,
				"balance": 0,
			},
			"withdrawn": cardBalance,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func PullGpuHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gpuIdStr := chi.URLParam(r, "gpuId")
		userId := chi.URLParam(r, "userId")

		gpuId, err := strconv.Atoi(gpuIdStr)
		if err != nil {
			http.Error(w, "Invalid gpuId", http.StatusBadRequest)
			return
		}

		user, err := database.GetUser(db, userId)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		card, err := database.GetCardById(db, gpuId)
		if err != nil {
			http.Error(w, "Card not found", http.StatusNotFound)
			return
		}

		if card.UserId != user.Id {
			json.NewEncoder(w).Encode(map[string]string{"status": "dontHaveGpu"})
			return
		}

		stand, err := database.GetStandByCardId(db, card.Id)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"status": "dontHaveStand"})
			return
		}

		if err := database.RemoveCardFromStand(db, stand.Id); err != nil {
			http.Error(w, "Failed to remove card from stand", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"status": "success",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
