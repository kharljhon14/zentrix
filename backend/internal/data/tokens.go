package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kharljhon14/zentrix/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    uuid.UUID `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

func (t TokenModel) New(userID uuid.UUID, ttl time.Duration, scope string) (*Token, error) {
	token := generateToken(userID, ttl, scope)

	err := t.Insert(token)

	return token, err
}

func (t TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens 
		(hash, user_id, expiry, scope)
		VALUES
		($1, $2, $3, $4)
	`
	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.DB.ExecContext(ctx, query, args...)

	return err
}

func (t TokenModel) GetForToken(tokenScope, plainTextToken string) (*User, error) {
	hashedToken := sha256.Sum256([]byte(plainTextToken))

	query := `
		SELECT u.id, u.first_name, u.last_name, u.email, u.role, u.created_at, u.updated_at
		FROM users u
		JOIN tokens t
		ON u.id = t.user_id
		WHERE t.hash = $1
		AND t.scope = $2
		AND t.expiry > NOW()
	`

	args := []any{hashedToken[:], tokenScope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, sql.ErrNoRows
		default:
			return nil, err
		}
	}

	return &user, nil

}

func (t TokenModel) DeleteAllForUser(scope string, userID uuid.UUID) error {
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.DB.ExecContext(ctx, query, scope, userID)

	return err

}

func ValidatePlainTextToken(v *validator.Validator, plainTextToken string) {
	v.Check(plainTextToken != "", "token", "token is required")
	v.Check(len(plainTextToken) == 26, "token", "must be 26 bytes long")
}

func generateToken(userID uuid.UUID, ttl time.Duration, scope string) *Token {
	token := &Token{
		PlainText: rand.Text(),
		UserID:    userID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}

	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token
}
