package domain

import "time"

type Permission struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PermissionRelation struct {
	Permission
}
