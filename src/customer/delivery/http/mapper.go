package http

import (
	"github.com/rasteiro11/MCABankGateway/src/customer/domain"
)

func MapCreateRequestToDomain(req *createCustomerRequest) *domain.Customer {
	return &domain.Customer{
		Nome:  req.Nome,
		Email: req.Email,
	}
}

func MapUpdateRequestToDomain(id uint, req *updateCustomerRequest) *domain.Customer {
	return &domain.Customer{
		ID:    id,
		Nome:  req.Nome,
		Email: req.Email,
	}
}

func MapCustomerToResponse(c *domain.Customer) *customerResponse {
	return &customerResponse{
		ID:    c.ID,
		Nome:  c.Nome,
		Email: c.Email,
		Saldo: c.Saldo,
	}
}

func MapCustomersToResponse(customers []*domain.Customer) []*customerResponse {
	responses := make([]*customerResponse, len(customers))
	for i, c := range customers {
		responses[i] = MapCustomerToResponse(c)
	}
	return responses
}
