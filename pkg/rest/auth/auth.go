package auth

import (
	"context"
	"net/http"

	client "github.com/rasteiro11/MCABankGateway/pkg/transport/http"
	"github.com/rasteiro11/MCABankGateway/src/auth/domain"
	"github.com/rasteiro11/PogCore/pkg/logger"
)

var _ AuthClient = (*authClient)(nil)

type authClient struct {
	BaseURL string
}

func New(baseURL string) AuthClient {
	return &authClient{BaseURL: baseURL}
}

func (a *authClient) Login(ctx context.Context, email, password string) (*domain.AuthResponse, error) {
	path := a.BaseURL + "/auth/signin"
	req := &loginRequest{
		Email:    email,
		Password: password,
	}
	res, err := client.Post[domain.AuthResponse](ctx, path, req)
	if err != nil {
		if httpErr, ok := err.(*client.HTTPError); ok {
			if httpErr.StatusCode == http.StatusUnauthorized {
				return nil, domain.ErrInvalidCredentials
			}
		}
		return nil, err
	}
	return res, nil
}

func (a *authClient) Register(ctx context.Context, email, password string) (*domain.AuthResponse, error) {
	path := a.BaseURL + "/auth/register"
	req := &registerRequest{
		Email:    email,
		Password: password,
	}
	res, err := client.Post[domain.AuthResponse](ctx, path, req)
	if err != nil {
		if httpErr, ok := err.(*client.HTTPError); ok {
			logger.Of(ctx).Errorf("request failed: %s", httpErr.Error())
			return nil, httpErr
		}
		return nil, err
	}

	return res, nil
}
