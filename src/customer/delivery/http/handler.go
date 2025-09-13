package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rasteiro11/MCABankGateway/src/customer/service"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/server"
	"github.com/rasteiro11/PogCore/pkg/transport/rest"
)

var CustomerGroupPath = "/clientes"

type (
	HandlerOpt func(*handler)
	handler    struct {
		customerService service.CustomerService
	}
)

func WithCustomerService(cs service.CustomerService) HandlerOpt {
	return func(h *handler) {
		h.customerService = cs
	}
}

func NewHandler(server server.Server, opts ...HandlerOpt) {
	h := &handler{}
	for _, opt := range opts {
		opt(h)
	}

	server.AddHandler("", CustomerGroupPath, http.MethodGet, h.GetAll)
	server.AddHandler("/:id", CustomerGroupPath, http.MethodGet, h.GetByID)
	server.AddHandler("", CustomerGroupPath, http.MethodPost, h.Create)
	server.AddHandler("/:id", CustomerGroupPath, http.MethodPut, h.Update)
	server.AddHandler("/:id", CustomerGroupPath, http.MethodDelete, h.Delete)
}

var ErrInvalidRequest = errors.New("invalid request")

func (h *handler) GetAll(c *fiber.Ctx) error {
	customers, err := h.customerService.GetAll(c.Context())
	if err != nil {
		return rest.NewStatusInternalServerError(c, err)
	}

	jsonData, _ := json.Marshal(customers)

	logger.Of(c.Context()).Infof("retrieved %s", string(jsonData))
	resp := MapCustomersToResponse(customers)
	return rest.NewStatusOk(c, rest.WithBody(resp))
}

func (h *handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return rest.NewStatusBadRequest(c, ErrInvalidRequest)
	}

	customer, err := h.customerService.GetByID(c.Context(), uint(id))
	if err != nil {
		return rest.NewStatusInternalServerError(c, err)
	}

	resp := MapCustomerToResponse(customer)
	return rest.NewStatusOk(c, rest.WithBody(resp))
}

func (h *handler) Create(c *fiber.Ctx) error {
	req := &createCustomerRequest{}
	if err := c.BodyParser(req); err != nil {
		return rest.NewStatusBadRequest(c, ErrInvalidRequest)
	}

	customer := MapCreateRequestToDomain(req)
	created, err := h.customerService.Create(c.Context(), customer)
	if err != nil {
		return rest.NewStatusInternalServerError(c, err)
	}

	resp := MapCustomerToResponse(created)
	return rest.NewStatusCreated(c, rest.WithBody(resp))
}

func (h *handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return rest.NewStatusBadRequest(c, ErrInvalidRequest)
	}

	req := &updateCustomerRequest{}
	if err := c.BodyParser(req); err != nil {
		return rest.NewStatusBadRequest(c, ErrInvalidRequest)
	}

	customer := MapUpdateRequestToDomain(uint(id), req)
	updated, err := h.customerService.Update(c.Context(), customer)
	if err != nil {
		return rest.NewStatusInternalServerError(c, err)
	}

	resp := MapCustomerToResponse(updated)
	return rest.NewStatusOk(c, rest.WithBody(resp))
}

func (h *handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return rest.NewStatusBadRequest(c, ErrInvalidRequest)
	}

	if err := h.customerService.Delete(c.Context(), uint(id)); err != nil {
		return rest.NewStatusInternalServerError(c, err)
	}

	return rest.NewStatusOk(c, rest.WithBody(fiber.Map{"status": "deleted"}))
}
