package port

import (
	"time"

	"github.com/google/uuid"
)

type EmailLog struct {
	ID             uuid.UUID      `json:"id"`
	TemplateID     uuid.UUID      `json:"templateId"`
	WorkspaceID    uuid.UUID      `json:"workspaceId"`
	RecipientEmail string         `json:"recipientEmail"`
	Event          string         `json:"event"`
	Subject        string         `json:"subject"`
	Content        string         `json:"content"`
	Status         EmailLogStatus `json:"status"`
	ErrorMessage   string         `json:"errorMessage"`
	SentAt         time.Time      `json:"sentAt"`
}
