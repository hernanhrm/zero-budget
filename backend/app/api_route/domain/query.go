package domain

import "github.com/google/uuid"

type ApiRoute struct {
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	PermissionID uuid.UUID `json:"permissionId"`
}

type ApiRouteRelation struct {
	ApiRoute
}
