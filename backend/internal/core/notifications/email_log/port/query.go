package port

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type EmailLog struct {
	ID             uuid.UUID      `json:"id"`
	TemplateID     uuid.UUID      `json:"templateId"`
	OrganizationID null.String    `json:"organizationId"`
	RecipientEmail string         `json:"recipientEmail"`
	Event          string         `json:"event"`
	Subject        string         `json:"subject"`
	Content        string         `json:"content"`
	Status         EmailLogStatus `json:"status"`
	ErrorMessage   string         `json:"errorMessage"`
	SentAt         time.Time      `json:"sentAt"`
}
