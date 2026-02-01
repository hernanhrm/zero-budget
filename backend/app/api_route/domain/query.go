package domain

import (
	"time"

	"github.com/google/uuid"
)

type ApiRoute struct {
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	PermissionID uuid.UUID `json:"permissionId"`
	CreatedAt    time.Time `json:"createdAt"`
}

type ApiRouteRelation struct {
	ApiRoute
}
