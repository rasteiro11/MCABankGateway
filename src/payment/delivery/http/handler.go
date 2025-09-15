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

// Deposit godoc
// @Summary Deposit money
// @Description Deposit money to a customer's account
// @Tags Payment
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param request body paymentRequest true "Deposit request"
// @Success 200 {object} depositResponse
// @Failure 400 {object} any
// @Failure 500 {object} any
// @Router /clientes/{id}/depositar [post]
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

// Withdraw godoc
// @Summary Withdraw money
// @Description Withdraw money from a customer's account
// @Tags Payment
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param request body paymentRequest true "Withdraw request"
// @Success 200 {object} withdrawResponse
// @Failure 400 {object} any
// @Failure 500 {object} any
// @Router /clientes/{id}/sacar [post]
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
