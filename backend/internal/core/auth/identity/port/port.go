package port

import (
	"context"

	basedomain "backend/port"
)

type Service interface {
	basedomain.UseCaseTx[Service]
	SignupWithEmail(ctx context.Context, input SignupWithEmail) (SignupResponse, error)
	LoginWithEmail(ctx context.Context, input LoginWithEmail) (LoginResponse, error)
	RefreshToken(ctx context.Context, accessToken, refreshToken string) (RefreshResponse, error)
}
