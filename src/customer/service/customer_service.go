package service

import (
	"context"

	customerClient "github.com/rasteiro11/MCABankGateway/pkg/rest/customer"
	balanceService "github.com/rasteiro11/MCABankGateway/src/balance/service"
	"github.com/rasteiro11/MCABankGateway/src/customer/domain"
	"github.com/rasteiro11/PogCore/pkg/logger"
)

type customerService struct {
	customerClient customerClient.CustomerClient
	balanceService balanceService.BalanceService
}

var _ CustomerService = (*customerService)(nil)

func NewCustomerService(cc customerClient.CustomerClient, ps balanceService.BalanceService) CustomerService {
	return &customerService{
		customerClient: cc,
		balanceService: ps,
	}
}

func (s *customerService) GetAll(ctx context.Context) ([]*domain.Customer, error) {
	customers, err := s.customerClient.GetAll(ctx, nil)
	if err != nil {
		logger.Of(ctx).Error("failed to get all customers")
		return nil, err
	}

	logger.Of(ctx).Info("filling customer balances")
	if err := s.balanceService.FillCustomerBalances(ctx, customers); err != nil {
		logger.Of(ctx).Error("failed to fill customer balances")
		return nil, err
	}

	return customers, nil
}

func (s *customerService) GetByID(ctx context.Context, id uint) (*domain.Customer, error) {
	customer, err := s.customerClient.GetByID(ctx, id)
	if err != nil {
		logger.Of(ctx).Errorf("failed to get customer with id %d", id)
		return nil, err
	}
	return customer, nil
}

func (s *customerService) Create(ctx context.Context, c *domain.Customer) (*domain.Customer, error) {
	customer, err := s.customerClient.Create(ctx, c)
	if err != nil {
		logger.Of(ctx).Errorf("failed to create customer with id %d", c.ID)
		return nil, err
	}

	return customer, nil
}

func (s *customerService) Update(ctx context.Context, c *domain.Customer) (*domain.Customer, error) {
	customer, err := s.customerClient.Update(ctx, c.ID, c)
	if err != nil {
		logger.Of(ctx).Errorf("failed to update customer with id %d", c.ID)
		return nil, err
	}

	if err := s.balanceService.FillCustomerBalance(ctx, customer); err != nil {
		logger.Of(ctx).Errorf("failed to fill customer balance with id %d", c.ID)
		return nil, err
	}

	return customer, nil
}

func (s *customerService) Delete(ctx context.Context, id uint) error {
	err := s.customerClient.Delete(ctx, id)
	if err != nil {
		logger.Of(ctx).Errorf("failed to delete customer with id %d", id)
		return err
	}
	return nil
}
