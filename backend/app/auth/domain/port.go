package domain

import (
	"context"

	basedomain "backend/domain"
)

type Service interface {
	basedomain.UseCaseTx[Service]
	SignupWithEmail(ctx context.Context, input SignupWithEmail) (SignupResponse, error)
	LoginWithEmail(ctx context.Context, input LoginWithEmail) (LoginResponse, error)
	RefreshToken(ctx context.Context, accessToken, refreshToken string) (RefreshResponse, error)
}
