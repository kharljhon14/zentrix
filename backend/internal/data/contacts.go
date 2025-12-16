package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
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

type ContactModel struct {
	DB *sql.DB
}

func (c ContactModel) Insert(contact Contact) error {
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
		WHERE id = $1 AND deleted_at IS NULL
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var contact ContactWithCompanyName
	err := c.DB.QueryRowContext(ctx, query, ID).Scan(
		&contact.ID,
		&contact.Name,
		&contact.Email,
		&contact.CompanyID,
		&contact.CompanyID,
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
