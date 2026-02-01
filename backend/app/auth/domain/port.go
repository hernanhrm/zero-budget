package domain

import (
	"context"
)

type Service interface {
	SignupWithEmail(ctx context.Context, input SignupWithEmail) (SignupResponse, error)
}
