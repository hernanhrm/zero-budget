package port

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type Category struct {
	ID             uuid.UUID   `json:"id"`
	OrganizationID string      `json:"organizationId"`
	ParentID       *uuid.UUID  `json:"parentId"`
	Name           string      `json:"name"`
	Icon           null.String `json:"icon"`
	Color          null.String `json:"color"`
	IsActive       bool        `json:"isActive"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}
