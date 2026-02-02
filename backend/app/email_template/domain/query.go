package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmailTemplate struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspaceId"`
	Event       string    `json:"event"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Subject     string    `json:"subject"`
	Content     string    `json:"content"`
	IsActive    bool      `json:"isActive"`
	Locale      string    `json:"locale"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
