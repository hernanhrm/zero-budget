package domain

import "github.com/guregu/null/v6"

type CreateUser struct {
	Name  string
	Email string
}

type UpdateUser struct {
	Name  null.String
	Email null.String
}
