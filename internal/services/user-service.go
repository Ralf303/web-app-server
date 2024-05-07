package user

import (
	"example.com/myapp/internal/database"
	"github.com/jmoiron/sqlx"
)

func UpdatedBalance(db *sqlx.DB, chatId string, value uint64) error {
	user, err := database.GetOrCreateUser(db, chatId)
	if err != nil {
		return err
	}
	err = database.UpdateUser(db, user.Id, "balance", value)
	if err != nil {
		return err
	}
	return nil
}

func UpdatedGems(db *sqlx.DB, chatId string, value uint64) error {
	user, err := database.GetOrCreateUser(db, chatId)
	if err != nil {
		return err
	}
	err = database.UpdateUser(db, user.Id, "gems", value)
	if err != nil {
		return err
	}
	return nil
}

func UpdatedKeys(db *sqlx.DB, chatId string, value uint64) error {
	user, err := database.GetOrCreateUser(db, chatId)
	if err != nil {
		return err
	}
	err = database.UpdateUser(db, user.Id, "chests", value)
	if err != nil {
		return err
	}
	return nil
}
