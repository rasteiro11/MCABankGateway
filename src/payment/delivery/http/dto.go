package http

type paymentRequest struct {
	CustomerId     uint    `json:"customer_id"`
	Amount         float64 `json:"amount"`
	IdempotencyKey string  `json:"idempotency_key"`
}

type depositResponse struct {
	Status string `json:"status"`
}

type withdrawResponse struct {
	Status string `json:"status"`
}
