package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/kharljhon14/zentrix/internal/validator"
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

func (p ProductModel) Insert(product *Product) error {
	query := `
		INSERT INTO products
			(quote_id, title, unit_price, quantity, discount)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		product.QuoteID,
		product.Title,
		product.UnitPrice,
		product.Quantity,
		product.Discount,
	}

	return p.DB.QueryRowContext(ctx, query, args...).Scan(
		&product.ID,
	)
}

func (p ProductModel) GetByQuoteID(ID uuid.UUID) ([]*Product, error) {
	return nil, nil
}

func (p Product) ValidateProduct(v *validator.Validator) {
	v.Check(p.Title != "", "title", "title is required")
	v.Check(len(p.Title) < 255, "title", "title must not exceed 255 chaaracters")
	v.Check(p.UnitPrice > 0, "unit_price", "unit_price mut be valid")
	v.Check(p.UnitPrice < 10_000_000, "unit_price", "unit_price must not exceed 10,000,000")
	v.Check(p.Quantity < 1_000_000, "quantity", "quantity must not exceed 1,000,000")
	v.Check(p.Discount < 100, "quantity", "quantity must not exceed 100")
}
