// Package errors provides shared error codes and layer constants for consistent error handling.
package errors

// Error codes represent expected error conditions mapped to HTTP status codes.
// Use these with oops.Code() only for expected errors (validation, not found, etc.).
// Unexpected errors should not have a code.
const (
	// CodeNotFound indicates a requested resource was not found (404).
	CodeNotFound = "not_found"

	// CodeBadRequest indicates malformed request syntax or invalid parameters (400).
	CodeBadRequest = "bad_request"

	// CodeValidation indicates request validation failed (422).
	CodeValidation = "validation"

	// CodeConflict indicates a conflict with current resource state (409).
	CodeConflict = "conflict"

	// CodeUnauthorized indicates missing or invalid authentication (401).
	CodeUnauthorized = "unauthorized"

	// CodeForbidden indicates insufficient permissions (403).
	CodeForbidden = "forbidden"

	// CodeAlreadyExists indicates a resource already exists (409).
	CodeAlreadyExists = "already_exists"
)
