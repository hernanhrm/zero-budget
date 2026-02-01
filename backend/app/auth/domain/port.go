package domain

import (
	"context"

	basedomain "backend/domain"
)

type Service interface {
	basedomain.UseCaseTx[Service]
	SignupWithEmail(ctx context.Context, input SignupWithEmail) (SignupResponse, error)
}
