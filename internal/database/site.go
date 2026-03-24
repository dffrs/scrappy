package database

import (
	"context"
	"database/sql"
	"time"
)

type SiteModel struct {
	DB *sql.DB
}

type Site struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (s *SiteModel) GetOrCreate(name string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	insert := `INSERT OR IGNORE INTO site (name) VALUES (?)`

	_, err := s.DB.ExecContext(ctx, insert, name)
	if err != nil {
		return 0, err
	}

	query := "SELECT id FROM site WHERE name = ?"

	var id int
	err = s.DB.QueryRowContext(ctx, query, name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
