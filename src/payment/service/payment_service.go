package service

import (
	"context"

	pbPaymentClient "github.com/rasteiro11/MCABankGateway/gen/proto/go/payment"
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
	_, err := s.paymentClient.Deposit(ctx, &pbPaymentClient.DepositRequest{
		CustomerId:     uint32(customerID),
		Amount:         amount,
		IdempotencyKey: idempotencyKey,
	})
	return err
}

func (s *service) Transfer(ctx context.Context, customerID uint, amount float64, idempotencyKey string) error {
	_, err := s.paymentClient.Transfer(ctx, &pbPaymentClient.TransferRequest{
		CustomerId:     uint32(customerID),
		Amount:         amount,
		IdempotencyKey: idempotencyKey,
	})
	return err
}
