package database

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type Card struct {
	Id      int       `db:"id"`
	UserId  int       `db:"userId"`
	Lvl     int       `db:"lvl"`
	Fuel    int       `db:"fuel"`
	Balance int       `db:"balance"`
	Created time.Time `db:"createdAt"`
	Updated time.Time `db:"updatedAt"`
}

type CardStand struct {
	Id        int       `db:"id"`
	UserId    int       `db:"userId"`
	CardId    *int      `db:"cardId"`
	Card      *Card     `db:"card"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}

type CardStandJoined struct {
	Id        int       `db:"id"`
	UserId    int       `db:"userId"`
	CardId    *int      `db:"cardId"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`

	Card_Id        *int       `db:"card.id"`
	Card_UserId    *int       `db:"card.userId"`
	Card_Lvl       *int       `db:"card.lvl"`
	Card_Fuel      *int       `db:"card.fuel"`
	Card_Balance   *int       `db:"card.balance"`
	Card_CreatedAt *time.Time `db:"card.createdAt"`
	Card_UpdatedAt *time.Time `db:"card.updatedAt"`
}

func GetUserCardStands(db *sqlx.DB, userId int) ([]CardStand, error) {
	var rows []CardStandJoined
	err := db.Select(&rows, `
		SELECT 
			cs.id, 
			cs.userId, 
			cs.cardId, 
			cs.createdAt, 
			cs.updatedAt,
			c.id AS "card.id",
			c.userId AS "card.userId",
			c.lvl AS "card.lvl", 
			c.fuel AS "card.fuel", 
			c.balance AS "card.balance",
			c.createdAt AS "card.createdAt",
			c.updatedAt AS "card.updatedAt"
		FROM cardStands cs
		LEFT JOIN cards c ON cs.cardId = c.id
		WHERE cs.userId = ?`, userId)
	if err != nil {
		return nil, err
	}

	stands := make([]CardStand, 0, len(rows))
	for _, r := range rows {
		var card *Card
		if r.Card_Id != nil { // создаём карту только если она есть
			card = &Card{
				Id:      *r.Card_Id,
				UserId:  *r.Card_UserId,
				Lvl:     *r.Card_Lvl,
				Fuel:    *r.Card_Fuel,
				Balance: *r.Card_Balance,
				Created: *r.Card_CreatedAt,
				Updated: *r.Card_UpdatedAt,
			}
		}

		stands = append(stands, CardStand{
			Id:        r.Id,
			UserId:    r.UserId,
			CardId:    r.CardId,
			Card:      card,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}

	return stands, nil
}

func GetUserCards(db *sqlx.DB, userId int) ([]Card, error) {
	var cards []Card
	err := db.Select(&cards, "SELECT * FROM cards WHERE userId = ?", userId)
	return cards, err
}

func GetCardById(db *sqlx.DB, gpuId int) (Card, error) {
	var card Card
	err := db.Get(&card, "SELECT * FROM cards WHERE id = ?", gpuId)
	return card, err
}

func GetCardStandById(db *sqlx.DB, standId int) (CardStand, error) {
	var stand CardStand
	err := db.Get(&stand, "SELECT * FROM cardStands WHERE id = ?", standId)
	return stand, err
}

func ResetCardBalance(db *sqlx.DB, id int) error {
	_, err := db.Exec("UPDATE cards SET balance = 0, updatedAt = NOW() WHERE id = ?", id)
	return err
}

func CreateCardStand(db *sqlx.DB, userId int) (CardStand, error) {
	var stand CardStand
	_, err := db.Exec(`
		INSERT INTO cardStands (userId, cardId, createdAt, updatedAt)
		VALUES (?, NULL, NOW(), NOW())`, userId)
	if err != nil {
		return stand, err
	}
	err = db.Get(&stand, "SELECT * FROM cardStands WHERE userId = ? ORDER BY id DESC LIMIT 1", userId)
	return stand, err
}

func InsertCardIntoStand(db *sqlx.DB, standId int, cardId int) error {
	_, err := db.Exec(`
		UPDATE cardStands 
		SET cardId = ?, updatedAt = NOW() 
		WHERE id = ?`, cardId, standId)
	return err
}

func UpdateCardFuel(db *sqlx.DB, cardId int, fuel int) error {
	_, err := db.Exec(`
		UPDATE cards 
		SET fuel = ?, updatedAt = NOW() 
		WHERE id = ?`, fuel, cardId)
	return err
}

func RemoveCardFromStand(db *sqlx.DB, standId int) error {
	_, err := db.Exec(`
		UPDATE cardStands 
		SET cardId = NULL, updatedAt = NOW() 
		WHERE id = ?`, standId)
	return err
}

func IsCardInstalledElsewhere(db *sqlx.DB, cardId int) (bool, error) {
	var existingStandId int
	err := db.Get(&existingStandId, "SELECT id FROM cardStands WHERE cardId = ?", cardId)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
