package data

import (
	"database/sql"

	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID `json:"id"`
	QuoteID   uuid.UUID `json:"quote_id"`
	Title     string    `json:"title"`
	UnitPrice int       `json:"unit_price"`
	Quantity  int       `json:"quantity"`
	Discount  int       `json:"discount"`
}

type ProductModel struct {
	DB *sql.DB
}
