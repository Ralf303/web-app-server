package database

import (
	"database/sql"
	"fmt"
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
}

func GetAllUsers(db *sqlx.DB) []User {
	users := []User{}
	err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
	}
	return users
}

func GetOrCreateUser(db *sqlx.DB, chatId string) (User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE chatId = ?", chatId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Создаем нового пользователя
			_, err := db.Exec("INSERT INTO users (chatId, createdAt, updatedAt) VALUES (?, datetime('now'), datetime('now'))", chatId)
			if err != nil {
				fmt.Println("Ошибка при создании нового пользователя:", err)
				return User{}, err
			}

			err = db.Get(&user, "SELECT * FROM users WHERE chatId = ?", chatId)
			if err != nil {
				fmt.Println("Ошибка при получении информации о новом пользователе:", err)
				return User{}, err
			}

		} else {
			fmt.Println("Ошибка при поиске пользователя:", err)
			return User{}, err
		}
	}
	return user, nil
}

func UpdateUser(db *sqlx.DB, userId int, column string, value uint64) error {
	_, err := db.Exec("UPDATE users SET "+column+" = ? WHERE id = ?", value, userId)
	if err != nil {
		fmt.Println("Ошибка при обновлении пользователя:", err)
		return err
	}
	return nil
}
