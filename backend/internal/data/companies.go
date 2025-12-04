package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Company struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Email     string    `json:"email"`
	Image     *string   `json:"image"`
	Website   *string   `json:"website"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CompanyModel struct {
	DB *sql.DB
}

func (c CompanyModel) Insert(company *Company) error {
	query := `
		INSERT INTO companies 
		(name, address, email, image, website)
		VALUES 
		($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	args := []any{company.Name, company.Address, company.Email, company.Image, company.Website}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(
		&company.ID,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	if err != nil {
		fmt.Println(err)
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "companies_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return sql.ErrNoRows
		default:
			return err
		}
	}

	return nil
}
