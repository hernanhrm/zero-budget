package domain

import "time"

type Workspace struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organizationId"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type WorkspaceRelation struct {
	Workspace
}
