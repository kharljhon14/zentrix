package data

import (
	"context"
	"database/sql"
	"fmt"
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
	v.Check(len(project.Title) <= 255, "title", "title must not exceed 255 characters")
	v.Check(project.Description != "", "description", "description is required")
	v.Check(project.Status != "", "status", "status is required")
	v.Check(len(project.Status) <= 255, "status", "status must not exceed 255 characters")
}

func (p ProjectModel) GetByID(ID uuid.UUID) (*Project, error) {
	query := `
		SELECT 
			id,
			company_id,
			title,
			description,
			status,
			owner_id,
			created_at,
			updated_at
		FROM projects
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var project Project
	err := p.DB.QueryRowContext(ctx, query, ID).Scan(
		&project.ID,
		&project.CompanyID,
		&project.Title,
		&project.Description,
		&project.Status,
		&project.OwnerID,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (p ProjectModel) GetAllByCompanyID(ID uuid.UUID, filters Filters) ([]*Project, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT
			count(id) over(),
			id,
			company_id,
			title,
			description,
			status,
			owner_id,
			created_at,
			updated_at,
		FROM projects
		WHERE company_id = $1
		ORDER BY %s %s, c.created_at DESC
		LIMIT $2 OFFSET $3

	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{ID, filters.limit(), filters.offset()}

	rows, err := p.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	projects := []*Project{}

	for rows.Next() {
		var project Project

		err := rows.Scan(
			&totalRecords,
			&project.ID,
			&project.CompanyID,
			&project.Title,
			&project.Description,
			&project.Status,
			&project.OwnerID,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		projects = append(projects, &project)

	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return projects, metadata, nil
}
