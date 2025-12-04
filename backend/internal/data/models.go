package data

import (
	"database/sql"
	"errors"
)

var (
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	Users     UserModel
	Tokens    TokenModel
	Companies CompanyModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:     UserModel{DB: db},
		Tokens:    TokenModel{DB: db},
		Companies: CompanyModel{DB: db},
	}
}
