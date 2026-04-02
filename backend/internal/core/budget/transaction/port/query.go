package port

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v6"
)

type Transaction struct {
	ID                      uuid.UUID   `json:"id"`
	OrganizationID          string      `json:"organizationId"`
	AccountID               uuid.UUID   `json:"accountId"`
	CategoryID              *uuid.UUID  `json:"categoryId"`
	SubcategoryID           *uuid.UUID  `json:"subcategoryId"`
	BudgetID                *uuid.UUID  `json:"budgetId"`
	Type                    string      `json:"type"`
	Amount                  int64       `json:"amount"`
	Description             null.String `json:"description"`
	ExternalReferenceNumber null.String `json:"externalReferenceNumber"`
	Date                    time.Time   `json:"date"`
	CreatedAt               time.Time   `json:"createdAt"`
	UpdatedAt               time.Time   `json:"updatedAt"`
}
