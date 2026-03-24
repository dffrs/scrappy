package database

import (
	"context"
	"database/sql"
	"time"
)

type ProductModel struct {
	DB *sql.DB
}

type Product struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	SiteID *int   `json:"site_id"`
	URL    string `json:"url"`
}

func (s *ProductModel) GetOrCreate(name string, siteID int, url string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	insert := `INSERT OR IGNORE INTO product (name, site_id, url) VALUES (?, ?, ?)`

	_, err := s.DB.ExecContext(ctx, insert, name, siteID, url)
	if err != nil {
		return 0, err
	}

	query := "SELECT id FROM product WHERE name = ? AND site_id = ? AND url = ?"

	var id int
	err = s.DB.QueryRowContext(ctx, query, name, siteID, url).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
