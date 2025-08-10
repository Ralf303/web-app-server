package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"example.com/myapp/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type CaseOpenResponse struct {
	RewardType string `json:"reward_type"`
	Amount     uint64 `json:"amount"`
	KeysLeft   int    `json:"keys_left"`
}

func OpenCaseHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chatId := chi.URLParam(r, "chatId")
		user, err := database.GetOrCreateUser(db, chatId)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Chests < 1 {
			http.Error(w, "Not enough chests", http.StatusBadRequest)
			return
		}

		user.Chests -= 1
		err = database.UpdateUserKeys(db, user.Id, user.Chests)
		if err != nil {
			http.Error(w, "Failed to update chests", http.StatusInternalServerError)
			return
		}

		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		rewardRoll := rng.Intn(3)

		var rewardType string
		var amount uint64

		switch rewardRoll {
		case 0:
			rewardType = "gems"
			amount = uint64(rng.Intn(10) + 1)
			err = database.UpdateUserGems(db, user.Id, user.Gems+int(amount))
			if err != nil {
				http.Error(w, "Failed to update gems", http.StatusInternalServerError)
				return
			}
		case 1:
			rewardType = "balance"
			amount = uint64(rng.Intn(9001) + 1000)
			err = database.UpdateUserBalance(db, user.Id, user.Balance+amount)
			if err != nil {
				http.Error(w, "Failed to update balance", http.StatusInternalServerError)
				return
			}
		default:
			rewardType = "nothing"
			amount = 0
		}

		resp := CaseOpenResponse{
			RewardType: rewardType,
			Amount:     amount,
			KeysLeft:   user.Chests,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
