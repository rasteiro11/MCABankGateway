package service

import (
	"context"

	pbPaymentClient "github.com/rasteiro11/MCABankGateway/gen/proto/go/payment"
	"github.com/rasteiro11/PogCore/pkg/telemetry/tracer"
	"go.opentelemetry.io/otel/codes"
)

var _ PaymentService = (*service)(nil)

type service struct {
	paymentClient pbPaymentClient.PaymentServiceClient
}

func NewPaymentService(paymentClient pbPaymentClient.PaymentServiceClient) PaymentService {
	return &service{
		paymentClient: paymentClient,
	}
}

func (s *service) Deposit(ctx context.Context, customerID uint, amount float64, idempotencyKey string) error {
	ctx, span := tracer.Instance().Start(ctx, "payment.service.Deposit")
	defer span.End()

	span.SetAttribute("customer.id", int(customerID))
	span.SetAttribute("amount", amount)
	span.SetAttribute("idempotency.key", idempotencyKey)

	_, err := s.paymentClient.Deposit(ctx, &pbPaymentClient.DepositRequest{
		CustomerId:     uint32(customerID),
		Amount:         amount,
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		span.SetStatus(tracer.Status{Code: int64(codes.Error), Msg: ""})
		span.RecordError(err)
		return err
	}

	span.SetStatus(tracer.Status{Code: int64(codes.Ok), Msg: ""})
	return nil
}

func (s *service) Transfer(ctx context.Context, customerID uint, amount float64, idempotencyKey string) error {
	ctx, span := tracer.Instance().Start(ctx, "payment.service.Transfer")
	defer span.End()

	span.SetAttribute("customer.id", int(customerID))
	span.SetAttribute("amount", amount)
	span.SetAttribute("idempotency.key", idempotencyKey)

	_, err := s.paymentClient.Transfer(ctx, &pbPaymentClient.TransferRequest{
		CustomerId:     uint32(customerID),
		Amount:         amount,
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		span.SetStatus(tracer.Status{Code: int64(codes.Error), Msg: ""})
		span.RecordError(err)
		return err
	}

	span.SetStatus(tracer.Status{Code: int64(codes.Ok), Msg: ""})
	return nil
}
