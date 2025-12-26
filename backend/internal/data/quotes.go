package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Quote struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	CompanyID   uuid.UUID `json:"company_id"`
	TotalAmount int       `json:"total_amount"`
	SalesTax    int       `json:"sales_tax"`
	Stage       string    `json:"stage"`
	Notes       string    `json:"notes"`
	PreparedBy  uuid.UUID `json:"prepared_by"`
	PreparedFor uuid.UUID `json:"prepared_for"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type QuoteModel struct {
	DB *sql.DB
}

func (q QuoteModel) Insert(quote *Quote) error {
	query := `
		INSERT INTO quotes
			(name, company_id, total_amount, sales_tax, stage, notes, prepared_by, prepared_for)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	args := []any{
		quote.Name,
		quote.CompanyID,
		quote.TotalAmount,
		quote.SalesTax,
		quote.Stage,
		quote.Notes,
		quote.PreparedBy,
		quote.PreparedFor,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return q.DB.QueryRowContext(ctx, query, args...).Scan(
		&quote.ID,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	)
}

func (q QuoteModel) GetByID(ID uuid.UUID) (*Quote, error) {
	query := `
		SELECT 
			id,
			name,
			company_id,
			total_amount,
			sales_tax,
			stage,
			notes,
			prepared_by,
			prepared_for,
			created_at,
			updated_at
		FROM quotes
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var quote Quote
	err := q.DB.QueryRowContext(ctx, query, ID).Scan(
		&quote.ID,
		&quote.Name,
		&quote.CompanyID,
		&quote.SalesTax,
		&quote.Stage,
		&quote.Notes,
		&quote.PreparedBy,
		&quote.PreparedFor,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &quote, nil
}
