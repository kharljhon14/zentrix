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
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

func (p ProductModel) GetProductsByQuoteID(ID uuid.UUID) ([]*Product, error) {
	query := `
		SELECT
			id,
			quote_id,
			title,
			unit_price,
			quantity,
			discount
		FROM products
		WHERE quote_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, query, ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []*Product{}
	for rows.Next() {
		var product Product

		err := rows.Scan(
			&product.ID,
			&product.QuoteID,
			&product.Title,
			&product.UnitPrice,
			&product.Quantity,
			&product.Discount,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil

}

func (p ProductModel) Update(product *Product) (*Product, error) {
	query := `
		UPDATE products
		SET title = $1,
		unit_price = $2,
		quantity = $3,
		discount = $4,
		updated_at = now()
		WHERE id = $1
		RETURNING updated_at
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		product.Title,
		product.UnitPrice,
		product.Quantity,
		product.Discount,
	}

	err := p.DB.QueryRowContext(ctx, query, args...).Scan(
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil

}

func (p Product) ValidateProduct(v *validator.Validator) {
	v.Check(p.Title != "", "title", "title is required")
	v.Check(len(p.Title) < 255, "title", "title must not exceed 255 chaaracters")
	v.Check(p.UnitPrice > 0, "unit_price", "unit_price mut be valid")
	v.Check(p.UnitPrice < 10_000_000, "unit_price", "unit_price must not exceed 10,000,000")
	v.Check(p.Quantity < 1_000_000, "quantity", "quantity must not exceed 1,000,000")
	v.Check(p.Discount < 100, "quantity", "quantity must not exceed 100")
}
