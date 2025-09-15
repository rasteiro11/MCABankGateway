package http

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rasteiro11/MCABankGateway/src/auth/domain"
	"github.com/rasteiro11/MCABankGateway/src/auth/service"
	"github.com/rasteiro11/PogCore/pkg/server"
	"github.com/rasteiro11/PogCore/pkg/transport/rest"
)

var AuthGroupPath = "/auth"

type (
	HandlerOpt func(*handler)
	handler    struct {
		authService service.AuthService
	}
)

func WithAuthService(authService service.AuthService) HandlerOpt {
	return func(u *handler) {
		u.authService = authService
	}
}

func NewHandler(server server.Server, opts ...HandlerOpt) {
	h := &handler{}

	for _, opt := range opts {
		opt(h)
	}

	server.AddHandler("/signin", AuthGroupPath, http.MethodPost, h.Login)
	server.AddHandler("/register", AuthGroupPath, http.MethodPost, h.Register)
}

var ErrNotAuthorized = errors.New("not authorized")

var _ Handler = (*handler)(nil)

// Login godoc
// @Summary Sign in
// @Description Authenticate user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login request"
// @Success 200 {object} loginResponse
// @Failure 400 {object} any
// @Failure 401 {object} any
// @Failure 500 {object} any
// @Router /auth/signin [post]
func (h *handler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return rest.NewStatusBadRequest(c, err)
	}

	token, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return rest.NewStatusUnauthorized(c, ErrNotAuthorized)
		}
		return rest.NewStatusInternalServerError(c, err)
	}

	resp := loginResponse{Token: token.Token, ExpiresAt: token.ExpiresAt}
	return rest.NewStatusOk(c, rest.WithBody(resp))
}

// Register godoc
// @Summary Register
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Register request"
// @Success 201 {object} loginResponse
// @Failure 400 {object} any
// @Failure 401 {object} any
// @Failure 500 {object} any
// @Router /auth/register [post]
func (h *handler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return rest.NewStatusBadRequest(c, err)
	}

	token, err := h.authService.Register(c.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return rest.NewStatusUnauthorized(c, ErrNotAuthorized)
		}
		return rest.NewStatusInternalServerError(c, err)
	}

	resp := loginResponse{Token: token.Token, ExpiresAt: token.ExpiresAt}
	return rest.NewStatusCreated(c, rest.WithBody(resp))
}
