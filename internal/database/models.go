// Package database
package database

import "database/sql"

type Models struct {
	Site    SiteModel
	Product ProductModel
	History HistoryModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Site:    SiteModel{DB: db},
		Product: ProductModel{DB: db},
		History: HistoryModel{DB: db},
	}
}
