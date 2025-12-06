package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kharljhon14/zentrix/internal/validator"
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

func (c CompanyModel) GetByID(ID uuid.UUID) (*Company, error) {
	query := `
		SELECT id, name, address, email, image, created_at, updated_at
		FROM companies
		WHERE ID = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var company Company
	err := c.DB.QueryRowContext(ctx, query, ID).Scan(
		&company.ID,
		&company.Name,
		&company.Address,
		&company.Email,
		&company.Image,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (c CompanyModel) GetAll(filters Filters) ([]*Company, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(id) OVER(), id, name, address, email, image, website, created_at, updated_at
		FROM companies
		ORDER BY %s %s, id ASC
		LIMIT $1 OFFSET $2
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{filters.limit(), filters.offset()}

	rows, err := c.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	companies := []*Company{}

	for rows.Next() {
		var company Company

		err := rows.Scan(
			&totalRecords,
			&company.ID,
			&company.Name,
			&company.Address,
			&company.Email,
			&company.Image,
			&company.Website,
			&company.CreatedAt,
			&company.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		companies = append(companies, &company)
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return companies, metadata, nil
}

func (c CompanyModel) Update(company *Company) error {
	query := `
		UPDATE companies
		SET name = $1,
		address = $2,
		email = $3,
		image = $4,
		website = $5,
		updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at;
	`

	args := []any{
		company.Name,
		company.Address,
		company.Email,
		company.Image,
		company.Website,
		company.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(
		&company.UpdatedAt,
	)
	if err != nil {

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

func ValidateCompany(v *validator.Validator, company *Company) {
	v.Check(company.Name != "", "name", "name is required")
	v.Check(len(company.Name) <= 255, "name", "name must not exceed 255 characters")

	v.Check(company.Address != "", "address", "address is required")
	v.Check(len(company.Address) <= 255, "address", "address must not exceed 255 characters")

	v.Check(company.Email != "", "email", "email is required")

	if !validator.Matches(company.Email, validator.EmailRX) {
		v.AddError("email", "invalid email")
	}
}
