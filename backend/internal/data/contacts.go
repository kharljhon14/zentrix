package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kharljhon14/zentrix/internal/validator"
)

type Contact struct {
	ID        uuid.UUID  `json:"uuid"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CompanyID *uuid.UUID `json:"company_id"`
	Title     string     `json:"title"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (c Contact) ValidateContact(v *validator.Validator) {
	v.Check(c.Name != "", "name", "name is required")
	v.Check(len(c.Name) <= 255, "name", "name must not exceed 255 characters")

	v.Check(c.Email != "", "email", "email is required")
	if !validator.Matches(c.Email, validator.EmailRX) {
		v.AddError("email", "invalid email")
	}

	v.Check(c.Title != "", "title", "title is required")
	v.Check(len(c.Title) <= 255, "title", "title must not exceed 255 characters")
	v.Check(c.Status != "", "status", "status is required")
	v.Check(len(c.Status) <= 255, "status", "status must not exceed 255 characters")
}

type ContactModel struct {
	DB *sql.DB
}

func (c ContactModel) Insert(contact *Contact) error {
	query := `
		INSERT INTO contacts
		(name, email, company_id, title, status)
		VALUES
		($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	args := []any{
		contact.Name,
		contact.Email,
		contact.CompanyID,
		contact.Title,
		contact.Status,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(
		&contact.ID,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: insert or update on table "contacts" violates foreign key constraint "fk_company_id"`:
			return ErrInvalidUUID

		case err.Error() == `pq: duplicate key value violates unique constraint "contacts_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

type ContactWithCompanyName struct {
	ID          uuid.UUID  `json:"uuid"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	CompanyID   *uuid.UUID `json:"company_id"`
	CompanyName *string    `json:"company_name"`
	Title       string     `json:"title"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (c ContactModel) GetByIDWithCompanyName(ID uuid.UUID) (*ContactWithCompanyName, error) {
	query := `
		SELECT 
			c.id,
			c.name,
			c.email,
			o.id AS company_id,
			o.name AS company_name,
			c.title,
			c.status,
			c.created_at,
			c.updated_at
		FROM contacts c
		JOIN companies o
		ON c.company_id = o.id
		WHERE c.id = $1 AND c.deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var contact ContactWithCompanyName
	err := c.DB.QueryRowContext(ctx, query, ID).Scan(
		&contact.ID,
		&contact.Name,
		&contact.Email,
		&contact.CompanyID,
		&contact.CompanyName,
		&contact.Title,
		&contact.Status,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (c ContactModel) GetAll(filter Filters, companyID *uuid.UUID) ([]*ContactWithCompanyName, Metadata, error) {
	query := ""

	if companyID != nil {
		query = fmt.Sprintf(`
		SELECT 
			count(c.id) over(),
			c.id,
			c.name,
			c.email,
			o.id AS company_id,
			o.name AS company_name,
			c.title,
			c.status,
			c.created_at,
			c.updated_at
		FROM contacts c
		JOIN companies o
		ON c.company_id = o.id
		WHERE o.id = '%s' AND c.deleted_at IS NULL
		ORDER BY %s %s, c.created_at DESC
		LIMIT $1 OFFSET $2
		
	`, *companyID, filter.sortColumn(), filter.sortDirection())
	} else {
		query = fmt.Sprintf(`
		SELECT 
			count(c.id) over(),
			c.id,
			c.name,
			c.email,
			o.id AS company_id,
			o.name AS company_name,
			c.title,
			c.status,
			c.created_at,
			c.updated_at
		FROM contacts c
		JOIN companies o
		ON c.company_id = o.id
		WHERE c.deleted_at IS NULL
		ORDER BY %s %s, c.created_at DESC
		LIMIT $1 OFFSET $2
		
	`, filter.sortColumn(), filter.sortDirection())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{filter.limit(), filter.offset()}

	rows, err := c.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	contacts := []*ContactWithCompanyName{}

	for rows.Next() {
		var contact ContactWithCompanyName

		err := rows.Scan(
			&totalRecords,
			&contact.ID,
			&contact.Name,
			&contact.Email,
			&contact.CompanyID,
			&contact.CompanyName,
			&contact.Title,
			&contact.Status,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		contacts = append(contacts, &contact)
	}

	metadata := calculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return contacts, metadata, nil
}

func (c ContactModel) Update(contact Contact) error {
	query := `
		UPDATE contacts
			SET name = $1,
			email = $2,
			company_id = $3,
			title = $4,
			status = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	args := []any{
		contact.Name,
		contact.Email,
		contact.CompanyID,
		contact.Title,
		contact.Status,
		contact.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(
		&contact.UpdatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "contacts_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (c ContactModel) Delete(ID uuid.UUID) error {
	query := `
		UPDATE contacts
			SET deleted_at = NOW()
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
