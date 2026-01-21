package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/kharljhon14/zentrix/internal/validator"
)

type Project struct {
	ID          uuid.UUID `json:"id"`
	CompanyID   uuid.UUID `json:"company_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	OwnerID     uuid.UUID `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectModel struct {
	DB *sql.DB
}

func (p ProjectModel) Insert(project *Project) error {
	query := `
		INSERT into products
		(company_id, title, description, status, owner_id)
		VALUES 
		($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		project.CompanyID,
		project.Title,
		project.Description,
		project.Status,
		project.OwnerID,
	}

	return p.DB.QueryRowContext(ctx, query, args...).Scan(
		project.ID,
		project.CreatedAt,
		project.UpdatedAt,
	)
}

func (project Project) Validate(v *validator.Validator) {
	v.Check(project.Title != "", "title", "title is required")
	v.Check(len(project.Title) > 255, "title", "title must not exceed 255 characters")
	v.Check(project.Description != "", "description", "description is required")

}
