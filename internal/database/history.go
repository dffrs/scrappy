package database

import (
	"context"
	"database/sql"
	"time"
)

type HistoryModel struct {
	DB *sql.DB
}

type History struct {
	ID        int       `json:"id"`
	ProductID *int      `json:"product_id"`
	Price     float32   `json:"price"`
	Date      time.Time `json:"date"`
}

func (s *HistoryModel) GetOrCreate(productID int, price float32) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	insert := `INSERT INTO history (product_id, price) VALUES (?, ?)`

	_, err := s.DB.ExecContext(ctx, insert, productID)
	if err != nil {
		return 0, err
	}

	query := "SELECT id FROM history ORDER BY date DESC LIMIT 1 WHERE product_id = ?"

	var id int
	err = s.DB.QueryRowContext(ctx, query, productID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
