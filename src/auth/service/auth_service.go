package service

import (
	"context"

	authClient "github.com/rasteiro11/MCABankGateway/pkg/rest/auth"
	"github.com/rasteiro11/MCABankGateway/src/auth/domain"
	"github.com/rasteiro11/PogCore/pkg/telemetry/tracer"
	"go.opentelemetry.io/otel/codes"
	grpcCodes "google.golang.org/grpc/codes"
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
	ctx, span := tracer.Instance().Start(ctx, "auth.service.Login")
	defer span.End()

	span.SetAttribute("auth.email", email)

	token, err := s.authClient.Login(ctx, email, password)
	if err != nil {
		span.SetStatus(tracer.Status{Code: int64(codes.Error), Msg: ""})
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(tracer.Status{Code: int64(codes.Ok), Msg: ""})
	return token, nil
}

func (s *authService) Register(ctx context.Context, email, password string) (*domain.AuthResponse, error) {
	ctx, span := tracer.Instance().Start(ctx, "auth.service.Register")
	defer span.End()

	span.SetAttribute("auth.email", email)

	token, err := s.authClient.Register(ctx, email, password)
	if err != nil {
		if status.Code(err) == grpcCodes.AlreadyExists {
			span.SetStatus(tracer.Status{Code: int64(codes.Error), Msg: ""})
			span.RecordError(err)
			return nil, domain.ErrInvalidCredentials
		}

		span.SetStatus(tracer.Status{Code: int64(codes.Error), Msg: ""})
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(tracer.Status{Code: int64(codes.Ok), Msg: ""})
	return token, nil
}
