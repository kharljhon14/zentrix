package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kharljhon14/zentrix/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var AnonymousUser = &User{}

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) Insert(user *User) error {
	query := `
		INSERT INTO USERS (first_name, last_name, email, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	args := []any{
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password.hash,
		user.Role,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET first_name = $1,
		last_name = $2,
		email = $3,
		activated = $4,
		role = $5,
		updated_at = NOW()
		RETURNING updated_at
	`

	args := []any{user.FirstName, user.LastName, user.Email, user.Activated, user.Role}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return sql.ErrNoRows
		default:
			return err
		}
	}

	return nil
}

func (u UserModel) GetByID(ID uuid.UUID) (*User, error) {
	query := `
		SELECT 
			id,
			first_name,
			last_name,
			email,
			activated,
			role,
			created_at,
			updated_at
		FROM USERS
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := u.DB.QueryRowContext(ctx, query, ID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Activated,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil

}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "email is required")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "invalid email")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "password is required")
	// Todo: Make a regex for the password checking
	v.Check(len(password) >= 8, "password", "password must be atleast 8 characters")
	v.Check(len(password) <= 255, "password", "password must not exceed 255 characters")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.FirstName != "", "first_name", "first_name is required")
	v.Check(len(user.FirstName) <= 80, "first_name", "first_name must not exceed 80 characters")

	v.Check(user.LastName != "", "last_name", "last_name is required")
	v.Check(len(user.LastName) <= 80, "last_name", "last_name must not exceed 80 characters")

	v.Check(user.Role != "", "role", "role is required")
	v.Check(len(user.Role) <= 60, "role", "role must not exceed 40 characters")

	ValidateEmail(v, user.Email)

	if user.Password.plainText != nil {
		ValidatePassword(v, *user.Password.plainText)
	}

	if user.Password.hash == nil {
		panic("missing hash password from user")
	}
}

type password struct {
	plainText *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = hash

	return nil
}

func (p password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
