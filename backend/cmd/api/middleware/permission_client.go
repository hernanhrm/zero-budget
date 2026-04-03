package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
	Error         any  `json:"error"`
	HasPermission bool `json:"success"`
}

type permissionsRequest struct {
	Permissions map[string][]string `json:"permissions"`
}

// HasPermission checks whether the user (identified by forwarded headers) has the given permission.
func (pc PermissionClient) HasPermission(ctx context.Context, permission string, headers http.Header) (bool, error) {
	url := fmt.Sprintf("%s/api/auth/organization/has-permission", pc.identityURL)

	parts := strings.Split(permission, ":")
	if len(parts) != 2 {
		return false, oops.In(apperrors.LayerMiddleware).Code(apperrors.CodeBadRequest).Errorf("invalid permission format: %s", permission)
	}
	module := parts[0]
	action := parts[1]

	permReq := permissionsRequest{
		Permissions: map[string][]string{
			module: {action},
		},
	}
	body, err := json.Marshal(permReq)
	if err != nil {
		return false, oops.In(apperrors.LayerMiddleware).Wrapf(err, "failed to marshal permission request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return false, oops.In(apperrors.LayerMiddleware).Wrapf(err, "failed to create permission request")
	}

	req.Header.Set("Content-Type", "application/json")

	// Forward cookies, authorization, and origin headers for session identification
	if cookie := headers.Get("Cookie"); cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	if auth := headers.Get("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	if apiKey := headers.Get("X-API-Key"); apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}
	req.Header.Set("Origin", "http://localhost:8080")

	resp, err := pc.httpClient.Do(req)
	if err != nil {
		return false, oops.In(apperrors.LayerMiddleware).Wrapf(err, "failed to call identity service")
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[permission-client] request body sent: %s", body)
	log.Printf("[permission-client] response status: %d, body: %s", resp.StatusCode, respBody)

	if resp.StatusCode == http.StatusUnauthorized {
		return false, oops.In(apperrors.LayerMiddleware).Code(apperrors.CodeUnauthorized).Errorf("user is not authenticated")
	}

	if resp.StatusCode != http.StatusOK {
		return false, oops.In(apperrors.LayerMiddleware).Code(apperrors.CodeForbidden).Errorf("permission check failed with status %d", resp.StatusCode)
	}

	var result hasPermissionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return false, oops.In(apperrors.LayerMiddleware).Wrapf(err, "failed to decode permission response")
	}

	return result.HasPermission, nil
}
