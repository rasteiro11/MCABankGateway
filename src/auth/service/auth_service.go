package service

import (
	"context"

	authClient "github.com/rasteiro11/MCABankGateway/pkg/rest/auth"
	"github.com/rasteiro11/MCABankGateway/src/auth/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authService struct {
	authClient authClient.AuthClient
}

var _ AuthService = (*authService)(nil)

func NewAuthService(authClient authClient.AuthClient) AuthService {
	return &authService{
		authClient: authClient,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.AuthResponse, error) {
	token, err := s.authClient.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *authService) Register(ctx context.Context, email, password string) (*domain.AuthResponse, error) {
	token, err := s.authClient.Register(ctx, email, password)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	return token, nil
}
