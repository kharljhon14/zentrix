package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kharljhon14/zentrix/internal/validator"
)

type Quote struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	CompanyID   uuid.UUID `json:"company_id"`
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
			(name, company_id, sales_tax, stage, notes, prepared_by, prepared_for)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	args := []any{
		quote.Name,
		quote.CompanyID,
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

type QuoteWithRelationNames struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	CompanyID       uuid.UUID `json:"company_id"`
	CompanyName     string    `json:"company_name"`
	SalesTax        int       `json:"sales_tax"`
	Stage           string    `json:"stage"`
	Notes           string    `json:"notes"`
	PreparedBy      uuid.UUID `json:"prepared_by"`
	PreparedByName  string    `json:"prepared_by_name"`
	PreparedFor     uuid.UUID `json:"prepared_for"`
	PreparedForName string    `json:"prepared_for_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (q QuoteModel) GetAll(filter Filters) ([]*QuoteWithRelationNames, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT
			count(q.id) over(),
			q.id,
			q.name,
			q.company_id,
			c.name AS company_name,
			q.sales_tax,
			q.stage,
			q.notes,
			cn.id AS prepared_by,
			cn.first_name || ' ' || cn.last_name AS prepared_by_name,
			cnb.id AS prepared_for,
			cnb.name AS prepared_for_name,
			q.created_at,
			q.updated_at
		FROM quotes q
		JOIN companies c
			ON q.company_id = c.id
		JOIN users cn
			ON q.prepared_by = cn.id
		JOIN contacts cnb
			ON q.prepared_for = cnb.id
		ORDER BY %s %s, q.created_at DESC
		LIMIT $1 OFFSET $2
	`, filter.sortColumn(), filter.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{filter.limit(), filter.offset()}

	rows, err := q.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	quotes := []*QuoteWithRelationNames{}

	for rows.Next() {
		var quote QuoteWithRelationNames
		err := rows.Scan(
			&totalRecords,
			&quote.ID,
			&quote.Name,
			&quote.CompanyID,
			&quote.CompanyName,
			&quote.SalesTax,
			&quote.Stage,
			&quote.Notes,
			&quote.PreparedBy,
			&quote.PreparedByName,
			&quote.PreparedFor,
			&quote.PreparedForName,
			&quote.CreatedAt,
			&quote.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		quotes = append(quotes, &quote)
	}

	metadata := calculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return quotes, metadata, nil
}

func (q QuoteModel) Update(quote *Quote) error {
	query := `
		UPDATE quotes
		SET name = $1,
		company_id = $2,
		prepared_by = $3,
		prepared_for = $4,
		stage = $5,
		notes = $6,
		updated_at = NOW()
		WHERE id = $7
		returning updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		quote.Name,
		quote.CompanyID,
		quote.PreparedBy,
		quote.PreparedFor,
		quote.Stage,
		quote.Notes,
		quote.ID,
	}

	return q.DB.QueryRowContext(ctx, query, args...).Scan(
		&quote.UpdatedAt,
	)

}

func (q QuoteModel) Delete(ID uuid.UUID) error {
	query := `
		DELETE FROM quotes
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := q.DB.ExecContext(ctx, query, ID)
	if err != nil {
		return err
	}

	affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil

}

func (q Quote) ValidateQuote(v *validator.Validator) {
	v.Check(q.Name != "", "name", "name is required")
	v.Check(len(q.Name) < 255, "name", "name must not exceed 255 characters")
	v.Check(q.SalesTax > -1, "sales_tax", "sales_tax must be valid")
	v.Check(q.Stage != "", "stage", "stage is required")
	v.Check(len(q.Stage) < 255, "stage", "stage must not exceed 255 characters")
	v.Check(len(q.Notes) < 10000, "notes", "notes must not exceed 10,000 characters")
}
