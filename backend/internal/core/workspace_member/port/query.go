package port

import (
	"time"

	"github.com/google/uuid"
)

type WorkspaceMember struct {
	WorkspaceID uuid.UUID `json:"workspaceId"`
	UserID      uuid.UUID `json:"userId"`
	RoleID      uuid.UUID `json:"roleId"`
	CreatedAt   time.Time `json:"createdAt"`
}
