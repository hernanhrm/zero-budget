package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	apperrors "backend/port/errors"

	"github.com/samber/oops"
)

// PermissionClient calls the identity service to check user permissions.
type PermissionClient struct {
	identityURL string
	httpClient  *http.Client
}

// NewPermissionClient creates a new permission client for the given identity service URL.
func NewPermissionClient(identityURL string) *PermissionClient {
	return &PermissionClient{
		identityURL: identityURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type hasPermissionResponse struct {
	HasPermission bool `json:"hasPermission"`
}

// HasPermission checks whether the user (identified by forwarded headers) has the given permission.
func (pc PermissionClient) HasPermission(ctx context.Context, permission string, headers http.Header) (bool, error) {
	url := fmt.Sprintf("%s/api/auth/organization/has-permission?permission=%s", pc.identityURL, permission)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, oops.In(apperrors.LayerMiddleware).Wrapf(err, "failed to create permission request")
	}

	// Forward cookies and authorization header for session identification
	if cookie := headers.Get("Cookie"); cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	if auth := headers.Get("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return false, oops.In(apperrors.LayerMiddleware).Wrapf(err, "failed to call identity service")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return false, oops.In(apperrors.LayerMiddleware).Code(apperrors.CodeUnauthorized).Errorf("user is not authenticated")
	}

	if resp.StatusCode != http.StatusOK {
		return false, oops.In(apperrors.LayerMiddleware).Code(apperrors.CodeForbidden).Errorf("permission check failed with status %d", resp.StatusCode)
	}

	var result hasPermissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, oops.In(apperrors.LayerMiddleware).Wrapf(err, "failed to decode permission response")
	}

	return result.HasPermission, nil
}
