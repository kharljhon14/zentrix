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
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	SalesOwner   uuid.UUID `json:"sales_owner"`
	Email        string    `json:"email"`
	CompanySize  string    `json:"company_size"`
	Industry     string    `json:"industry"`
	BusinessType string    `json:"business_type"`
	Country      string    `json:"country"`
	Image        *string   `json:"image"`
	Website      *string   `json:"website"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CompanyModel struct {
	DB *sql.DB
}

func (c CompanyModel) Insert(company *Company) error {
	query := `
		INSERT INTO companies 
		(name, address, sales_owner, email, company_size, industry, business_type, country, image, website)
		VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	args := []any{
		company.Name,
		company.Address,
		company.SalesOwner,
		company.Email,
		company.CompanySize,
		company.Industry,
		company.BusinessType,
		company.Country,
		company.Image,
		company.Website,
	}

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

type CompanyWithSalesOwner struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Address        string     `json:"address"`
	SalesOwner     *uuid.UUID `json:"sales_owner"`
	SalesOwnerName *string    `json:"sales_owner_name"`
	Email          string     `json:"email"`
	CompanySize    string     `json:"company_size"`
	Industry       string     `json:"industry"`
	BusinessType   string     `json:"business_type"`
	Country        string     `json:"country"`
	Image          *string    `json:"image"`
	Website        *string    `json:"website"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (c CompanyModel) GetByID(ID uuid.UUID) (*Company, error) {
	query := `
		SELECT 
			id, 
			name, 
			address,
			sales_owner,
			email, 
			company_size, 
			business_type,
			industry,
			country, 
			image, 
			website,
			created_at, 
			updated_at
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
		&company.SalesOwner,
		&company.Email,
		&company.CompanySize,
		&company.BusinessType,
		&company.Industry,
		&company.Country,
		&company.Image,
		&company.Website,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (c CompanyModel) GetByIDWithSalesOwner(ID uuid.UUID) (*CompanyWithSalesOwner, error) {
	query := `
		SELECT 
			c.id, 
			c.name, 
			c.address,
			c.sales_owner,
			CONCAT(u.first_name, ' ', u.last_name) AS sales_owner_name,
			c.email, 
			c.company_size, 
			c.business_type,
			c.industry, 
			c.country, 
			c.image, 
			c.website,
			c.created_at, 
			c.updated_at
		FROM companies c
		JOIN users u
		ON c.sales_owner = u.id
		WHERE c.id = $1 AND c.deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var company CompanyWithSalesOwner
	err := c.DB.QueryRowContext(ctx, query, ID).Scan(
		&company.ID,
		&company.Name,
		&company.Address,
		&company.SalesOwner,
		&company.SalesOwnerName,
		&company.Email,
		&company.CompanySize,
		&company.BusinessType,
		&company.Industry,
		&company.Country,
		&company.Image,
		&company.Website,
		&company.CreatedAt,
		&company.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (c CompanyModel) GetAll(filters Filters) ([]*CompanyWithSalesOwner, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT 
			count(c.id) over(),
			c.id, 
			c.name, 
			c.address,
			c.sales_owner,
			CONCAT(u.first_name, ' ', u.last_name) AS sales_owner_name,
			c.email, 
			c.company_size, 
			c.business_type,
			c.industry, 
			c.country, 
			c.image, 
			c.website,
			c.created_at, 
			c.updated_at
		FROM companies c
		JOIN users u
		ON c.sales_owner = u.id
		WHERE c.deleted_at IS NULL
		ORDER BY %s %s, c.created_at DESC
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
	companies := []*CompanyWithSalesOwner{}

	for rows.Next() {
		var company CompanyWithSalesOwner

		err := rows.Scan(
			&totalRecords,
			&company.ID,
			&company.Name,
			&company.Address,
			&company.SalesOwner,
			&company.SalesOwnerName,
			&company.Email,
			&company.CompanySize,
			&company.BusinessType,
			&company.Industry,
			&company.Country,
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
		sales_owner = $3,
		email = $4,
		company_size = $5,
		business_type = $6,
		country = $7,
		image = $8,
		website = $9,
		updated_at = NOW()
		WHERE id = $10
		RETURNING updated_at;
	`

	args := []any{
		company.Name,
		company.Address,
		company.SalesOwner,
		company.Email,
		company.CompanySize,
		company.BusinessType,
		company.Country,
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

func (c CompanyModel) Delete(ID uuid.UUID) error {
	query := `
		UPDATE companies
		set deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.ExecContext(ctx, query, ID)
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

func ValidateCompany(v *validator.Validator, company *Company) {

	v.Check(company.Name != "", "name", "name is required")
	v.Check(len(company.Name) <= 255, "name", "name must not exceed 255 characters")

	v.Check(company.Address != "", "address", "address is required")
	v.Check(len(company.Address) <= 255, "address", "address must not exceed 255 characters")

	v.Check(company.CompanySize != "", "company_size", "company_size is required")
	v.Check(len(company.CompanySize) <= 255, "company_size", "company_size must not exceed 255 characters")

	v.Check(company.Industry != "", "industry", "industry is required")
	v.Check(len(company.Industry) <= 255, "industry", "industry must not exceed 255 characters")

	v.Check(company.BusinessType != "", "business_type", "business_type is required")
	v.Check(len(company.BusinessType) <= 255, "business_type", "business_type must not exceed 255 characters")

	v.Check(company.Country != "", "country", "country is required")
	v.Check(len(company.Country) <= 255, "country", "country must not exceed 255 characters")

	v.Check(company.Email != "", "email", "email is required")

	if !validator.Matches(company.Email, validator.EmailRX) {
		v.AddError("email", "invalid email")
	}
}
