package service

import (
	"context"

	"github.com/rasteiro11/MCABankGateway/src/customer/domain"
)

type BalanceService interface {
	FillCustomerBalances(ctx context.Context, customers []*domain.Customer) error
	FillCustomerBalance(ctx context.Context, customer *domain.Customer) error
}
