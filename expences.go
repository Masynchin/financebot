package main

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type ExpenceService struct {
	db *sqlx.DB
}

type Expence struct {
	Id        int64
	UserID    int64 `db:"user_id"`
	Category  string
	Amount    int
	Timestamp time.Time
}

type TotalExpence struct {
	UserID   int64 `db:"user_id"`
	Category string
	Amount   int
}

const createQuery = `
CREATE TABLE IF NOT EXISTS expences (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	category TEXT NOT NULL,
	amount INTEGER NOT NULL,
	timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

// Create expence service
func NewExpenceService(databaseURL string) (*ExpenceService, error) {
	db, err := sqlx.Open("sqlite", databaseURL)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createQuery); err != nil {
		return nil, err
	}

	return &ExpenceService{db: db}, nil
}

// Close expences service db connection
func (e *ExpenceService) Close() {
	e.db.Close()
}

// Add new expence
func (e *ExpenceService) InsertExpence(userID int64, category string, amount int) (int64, error) {
	query := "INSERT INTO expences(user_id, category, amount) VALUES($1, $2, $3)"
	res, err := e.db.Exec(query, userID, category, amount)
	if err != nil {
		return 0, err
	}

	expenceID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return expenceID, nil
}

// Get total expence amount with given category
func (e *ExpenceService) GetExpence(userID int64, category string) (int, error) {
	query := `
		SELECT sum(amount)
		FROM expences
		WHERE user_id = $1 AND category = $2
		GROUP BY category
	`
	res, err := e.db.Query(query, userID, category)
	if err != nil {
		return 0, err
	}

	var amount int
	for res.Next() {
		if err := res.Scan(&amount); err != nil {
			return 0, err
		}
	}

	return amount, nil
}

// Get expence by its own ID
func (e *ExpenceService) GetExpenceByID(expenceID int64) (*Expence, error) {
	exp := Expence{}
	if err := e.db.Get(&exp, "SELECT * FROM expences WHERE id = $1", expenceID); err != nil {
		return nil, err
	}

	return &exp, nil
}

// Get exact user expences
func (e *ExpenceService) GetUserExpences(userID int64) ([]Expence, error) {
	userExpences := []Expence{}
	query := `
		SELECT *
		FROM expences
		WHERE user_id = $1 AND timestamp BETWEEN $2 AND $3
	`
	monthStart, monthEnd := getMonthLimits()
	if err := e.db.Select(&userExpences, query, userID, monthStart, monthEnd); err != nil {
		return nil, err
	}

	return userExpences, nil
}

// Get user expences grouped by category
func (e *ExpenceService) GetUserExpencesGroupedByCategory(userID int64) ([]TotalExpence, error) {
	userExpencesGroupedByCategory := []TotalExpence{}
	query := `
		SELECT user_id, category, sum(amount) as amount
		FROM expences
		WHERE user_id = $1 AND timestamp BETWEEN $2 AND $3
		GROUP BY category
	`
	monthStart, monthEnd := getMonthLimits()
	err := e.db.Select(&userExpencesGroupedByCategory, query, userID, monthStart, monthEnd)
	if err != nil {
		return nil, err
	}

	return userExpencesGroupedByCategory, nil
}

func getMonthLimits() (monthStart time.Time, monthEnd time.Time) {
	now := time.Now()
	monthStart = time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.Local)
	monthEnd = time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local)

	return
}

// Delete expence
func (e *ExpenceService) DeleteExpence(expenceID int64) error {
	query := "DELETE FROM expences WHERE id = $1"
	_, err := e.db.Exec(query, expenceID)
	return err
}
