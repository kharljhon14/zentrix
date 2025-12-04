package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Plan      string    `json:"plan"`
	CreatedAt time.Time `json:"created_at"`
}

type TenantModel struct {
	DB *sql.DB
}

func (t TenantModel) Insert(tenant *Tenant) error {
	query := `
		INSERT INTO tenants (name, plan)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	args := []any{tenant.Name, tenant.Plan}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(
		&tenant.ID,
		&tenant.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (t TenantModel) Get(id string) (*Tenant, error) {
	query := `
		SELECT id, name, plan, created_at 
		FROM tenants 
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var tenant Tenant

	err := t.DB.QueryRowContext(ctx, query, id).Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Plan,
		&tenant.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}
