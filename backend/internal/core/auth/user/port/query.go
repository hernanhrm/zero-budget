package port

import (
	"time"

	"github.com/guregu/null/v6"
)

type User struct {
	ID           string      `json:"id"`
	FirstName    string      `json:"firstName"`
	LastName     string      `json:"lastName"`
	Email        string      `json:"email"`
	PasswordHash null.String `json:"-"`
	ImageURL     null.String `json:"imageUrl"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

type UserRelation struct {
	User
}
