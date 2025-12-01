package data

import "database/sql"

type Models struct {
	Tenants TenantModel
	Users   UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Tenants: TenantModel{DB: db},
		Users:   UserModel{DB: db},
	}
}
