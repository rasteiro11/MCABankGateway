package service

import "context"

type PaymentService interface {
	Deposit(ctx context.Context, customerID uint, amount float64, idempotencyKey string) error
	Transfer(ctx context.Context, customerID uint, amount float64, idempotencyKey string) error
}
