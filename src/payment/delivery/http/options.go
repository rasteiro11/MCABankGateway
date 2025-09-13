package http

import "github.com/rasteiro11/MCABankGateway/src/payment/service"

func WithPaymentService(paymentService service.PaymentService) HandlerOpt {
	return func(u *handler) {
		u.paymentService = paymentService
	}
}
