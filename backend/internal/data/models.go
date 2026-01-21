package data

import (
	"database/sql"
	"errors"
)

var (
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrInvalidUUID    = errors.New("invalid id")
)

type Models struct {
	Users     UserModel
	Tokens    TokenModel
	Companies CompanyModel
	Contacts  ContactModel
	Quotes    QuoteModel
	Products  ProductModel
	Projects  ProjectModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:     UserModel{DB: db},
		Tokens:    TokenModel{DB: db},
		Companies: CompanyModel{DB: db},
		Contacts:  ContactModel{DB: db},
		Quotes:    QuoteModel{DB: db},
		Products:  ProductModel{DB: db},
		Projects:  ProjectModel{DB: db},
	}
}
