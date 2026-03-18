package port

import (
	"time"

	roleport "backend/core/auth/role/port"
	userport "backend/core/auth/user/port"

	"github.com/google/uuid"
)

type WorkspaceMember struct {
	WorkspaceID uuid.UUID `json:"workspaceId"`
	UserID      uuid.UUID `json:"userId"`
	RoleID      uuid.UUID `json:"roleId"`
	CreatedAt   time.Time `json:"createdAt"`
}

type WorkspaceMemberRelation struct {
	WorkspaceMember
	User *userport.User `json:"user,omitempty"`
	Role *roleport.Role `json:"role,omitempty"`
}
