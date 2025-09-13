package customer

import (
	"context"

	"github.com/rasteiro11/MCABankGateway/src/customer/domain"
)

type CustomerClient interface {
	GetAll(ctx context.Context, queryParams map[string]string) ([]*domain.Customer, error)
	GetByID(ctx context.Context, id uint) (*domain.Customer, error)
	Create(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	Update(ctx context.Context, id uint, customer *domain.Customer) (*domain.Customer, error)
	Delete(ctx context.Context, id uint) error
}
