package service

import (
	"context"

	pbPaymentClient "github.com/rasteiro11/MCABankGateway/gen/proto/go/payment"
	"github.com/rasteiro11/MCABankGateway/src/customer/domain"
	"github.com/rasteiro11/PogCore/pkg/logger"
)

type balanceService struct {
	balanceClient pbPaymentClient.BalanceServiceClient
}

var _ BalanceService = (*balanceService)(nil)

func NewBalanceService(bc pbPaymentClient.BalanceServiceClient) BalanceService {
	return &balanceService{
		balanceClient: bc,
	}
}

func (s *balanceService) FillCustomerBalances(ctx context.Context, customers []*domain.Customer) error {
	customerIds := make([]uint32, len(customers))
	for i, customer := range customers {
		customerIds[i] = uint32(customer.ID)
	}

	balancesResp, err := s.balanceClient.GetBalances(ctx, &pbPaymentClient.GetBalancesRequest{
		CustomerIds: customerIds,
	})
	if err != nil {
		logger.Of(ctx).Error("failed to get customer balances")
		return err
	}

	balancesMap := make(map[uint32]float64)
	for _, balance := range balancesResp.Balances {
		balancesMap[balance.CustomerId] = balance.Balance
	}

	for _, c := range customers {
		if balance, ok := balancesMap[uint32(c.ID)]; ok {
			b := balance
			c.Saldo = b
		}
	}

	return nil
}

func (s *balanceService) FillCustomerBalance(ctx context.Context, customer *domain.Customer) error {
	balancesResp, err := s.balanceClient.GetBalances(ctx, &pbPaymentClient.GetBalancesRequest{
		CustomerIds: []uint32{uint32(customer.ID)},
	})
	if err != nil {
		logger.Of(ctx).Errorf("failed to get balance for customer with id %d", customer.ID)
		return err
	}

	if len(balancesResp.Balances) > 0 {
		b := balancesResp.Balances[0].Balance
		customer.Saldo = b
	}

	return nil
}
