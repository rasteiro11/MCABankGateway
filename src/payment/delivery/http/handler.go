package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rasteiro11/MCABankGateway/src/payment/service"
	"github.com/rasteiro11/PogCore/pkg/server"
	"github.com/rasteiro11/PogCore/pkg/transport/rest"
)

var PaymentGRoupPath = "/clientes"

var ErrInvalidRequest = errors.New("invalid request")

var _ Handler = (*handler)(nil)

type (
	HandlerOpt func(*handler)
	handler    struct {
		paymentService service.PaymentService
	}
)

func NewHandler(server server.Server, opts ...HandlerOpt) {
	h := &handler{}

	for _, opt := range opts {
		opt(h)
	}

	server.AddHandler("/:id/depositar", PaymentGRoupPath, http.MethodPost, h.Deposit)
	server.AddHandler("/:id/sacar", PaymentGRoupPath, http.MethodPost, h.Withdraw)
}

func (h *handler) Deposit(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return rest.NewStatusBadRequest(c, ErrInvalidRequest)
	}

	req := &paymentRequest{}
	if err := c.BodyParser(req); err != nil {
		return rest.NewStatusBadRequest(c, ErrInvalidRequest)
	}

	if err := h.paymentService.Deposit(c.Context(), uint(id), req.Amount, req.IdempotencyKey); err != nil {
		return rest.NewStatusInternalServerError(c, err)
	}

	return rest.NewStatusOk(c, rest.WithBody(&depositResponse{Status: "Depósito iniciado com sucesso"}))
}

func (h *handler) Withdraw(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return rest.NewStatusBadRequest(c, ErrInvalidRequest)
	}

	req := &paymentRequest{}
	if err := c.BodyParser(req); err != nil {
		return rest.NewStatusBadRequest(c, ErrInvalidRequest)
	}

	if err := h.paymentService.Transfer(c.Context(), uint(id), req.Amount, req.IdempotencyKey); err != nil {
		return rest.NewStatusInternalServerError(c, err)
	}

	return rest.NewStatusOk(c, rest.WithBody(&withdrawResponse{Status: "Saque iniciado com sucesso"}))
}
