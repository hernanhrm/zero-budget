package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          string    `json:"id"`
	WorkspaceID uuid.UUID `json:"workspaceId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
