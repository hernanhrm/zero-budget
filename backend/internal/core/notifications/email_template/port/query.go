package port

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type ParsedTemplate struct {
	Subject string
	Content string
}

type EmailTemplate struct {
	ID                     uuid.UUID   `json:"id"`
	OrganizationID         null.String `json:"organizationId"`
	Event                  string    `json:"event"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	Subject                string    `json:"subject"`
	Content                string    `json:"content"`
	IsActive               bool      `json:"isActive"`
	Locale                 string    `json:"locale"`
	IsOrganizationTemplate bool      `json:"isOrganizationTemplate"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}
