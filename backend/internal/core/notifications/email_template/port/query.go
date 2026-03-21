package port

import (
	"time"

	"github.com/google/uuid"
)

type ParsedTemplate struct {
	Subject string
	Content string
}

type EmailTemplate struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID string    `json:"organizationId"`
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
