package auth

import (
	"context"

	"github.com/rasteiro11/MCABankGateway/src/auth/domain"
)

type AuthClient interface {
	Login(ctx context.Context, email, password string) (*domain.AuthResponse, error)
	Register(ctx context.Context, email, password string) (*domain.AuthResponse, error)
}
