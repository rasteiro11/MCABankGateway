package service

import (
	"context"

	pbPaymentClient "github.com/rasteiro11/MCABankGateway/gen/proto/go/payment"
	"github.com/rasteiro11/MCABankGateway/src/customer/domain"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/telemetry/tracer"
	"go.opentelemetry.io/otel/codes"
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
	ctx, span := tracer.Instance().Start(ctx, "balance.service.FillCustomerBalances")
	defer span.End()

	customerIds := make([]uint32, len(customers))
	for i, customer := range customers {
		customerIds[i] = uint32(customer.ID)
	}
	span.SetAttribute("customers.count", len(customers))

	balancesResp, err := s.balanceClient.GetBalances(ctx, &pbPaymentClient.GetBalancesRequest{
		CustomerIds: customerIds,
	})
	if err != nil {
		span.SetStatus(tracer.Status{Code: int64(codes.Error), Msg: ""})
		span.RecordError(err)
		logger.Of(ctx).Error("failed to get customer balances")
		return err
	}

	balancesMap := make(map[uint32]float64)
	for _, balance := range balancesResp.Balances {
		balancesMap[balance.CustomerId] = balance.Balance
	}

	for _, c := range customers {
		if balance, ok := balancesMap[uint32(c.ID)]; ok {
			c.Saldo = balance
		}
	}

	span.SetStatus(tracer.Status{Code: int64(codes.Ok), Msg: ""})
	return nil
}

func (s *balanceService) FillCustomerBalance(ctx context.Context, customer *domain.Customer) error {
	ctx, span := tracer.Instance().Start(ctx, "balance.service.FillCustomerBalance")
	defer span.End()

	span.SetAttribute("customer.id", int(customer.ID))

	balancesResp, err := s.balanceClient.GetBalances(ctx, &pbPaymentClient.GetBalancesRequest{
		CustomerIds: []uint32{uint32(customer.ID)},
	})
	if err != nil {
		span.SetStatus(tracer.Status{Code: int64(codes.Error), Msg: ""})
		span.RecordError(err)
		logger.Of(ctx).Errorf("failed to get balance for customer with id %d", customer.ID)
		return err
	}

	if len(balancesResp.Balances) > 0 {
		customer.Saldo = balancesResp.Balances[0].Balance
	}

	span.SetStatus(tracer.Status{Code: int64(codes.Ok), Msg: ""})
	return nil
}
