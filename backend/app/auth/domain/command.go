package domain

import (
	"context"
	"time"

	"backend/infra/validation"
)

type SignupWithEmail struct {
	Email            string `json:"email"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Password         string `json:"password"`
	OrganizationName string `json:"organizationName"`
}

func (s SignupWithEmail) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &s,
		validation.Field(&s.Email, validation.Required, validation.IsEmail),
		validation.Field(&s.FirstName, validation.Required, validation.Length(2, 100)),
		validation.Field(&s.LastName, validation.Required, validation.Length(2, 100)),
		validation.Field(&s.Password, validation.Required, validation.Length(6, 255)),
		validation.Field(&s.OrganizationName, validation.Required, validation.Length(2, 255)),
	)
}

type TokenPair struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type SignupResponse struct {
	TokenPair
	User      UserInfo      `json:"user"`
	Workspace WorkspaceInfo `json:"workspace"`
}

type UserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type WorkspaceInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type JWTClaims struct {
	UserID      string `json:"user_id"`
	WorkspaceID string `json:"workspace_id"`
}

type LoginWithEmail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l LoginWithEmail) Validate(ctx context.Context) error {
	return validation.ValidateStruct(ctx, &l,
		validation.Field(&l.Email, validation.Required, validation.IsEmail),
		validation.Field(&l.Password, validation.Required),
	)
}

type LoginResponse struct {
	TokenPair
	User      UserInfo      `json:"user"`
	Workspace WorkspaceInfo `json:"workspace"`
}
