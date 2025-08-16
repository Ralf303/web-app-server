package database

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id             int       `db:"id"`
	ChatId         string    `db:"chatId"`
	Username       string    `db:"username"`
	FirstName      string    `db:"firstname"`
	CaptureCounter int       `db:"captureCounter"`
	Balance        uint64    `db:"balance"`
	Meflvl         int       `db:"meflvl"`
	Timelvl        int       `db:"timelvl"`
	Farmtime       int       `db:"farmtime"`
	CreatedAt      time.Time `db:"createdAt"`
	UpdatedAt      time.Time `db:"updatedAt"`
	Slots          int       `db:"slots"`
	FullSlots      int       `db:"fullSlots"`
	Gems           int       `db:"gems"`
	TakeBonus      int       `db:"takeBonus"`
	Chests         int       `db:"chests"`
	FamMoney       int       `db:"famMoney"`
	Stones         int       `db:"stones"`
	Snows          int       `db:"snows"`
	Freeze         int       `db:"freeze"`
	Oil            int       `db:"oil"`
	Donate         int       `db:"donate"`
	Coin           int       `db:"coin"`
}

func GetUser(db *sqlx.DB, chatId string) (User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE chatId = ?", chatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, err
	}
	return user, nil
}

func UpdateUserKeys(db *sqlx.DB, userId int, keys int) error {
	_, err := db.Exec("UPDATE users SET chests = ?, updatedAt = NOW() WHERE id = ?", keys, userId)
	return err
}

func UpdateUserBalance(db *sqlx.DB, userId int, balance uint64) error {
	_, err := db.Exec("UPDATE users SET balance = ?, updatedAt = NOW() WHERE id = ?", balance, userId)
	return err
}

func UpdateUserGems(db *sqlx.DB, userId int, gems int) error {
	_, err := db.Exec("UPDATE users SET gems = ?, updatedAt = NOW() WHERE id = ?", gems, userId)
	return err
}

func UpdateUserCoins(db *sqlx.DB, userId int, coins int) error {
	_, err := db.Exec("UPDATE users SET coin = ?, updatedAt = NOW() WHERE id = ?", coins, userId)
	return err
}

func DecrementUserFreeze(db *sqlx.DB, userId int) error {
	_, err := db.Exec(`
		UPDATE users 
		SET freeze = freeze - 1, updatedAt = NOW() 
		WHERE id = ? AND freeze > 0`, userId)
	return err
}
