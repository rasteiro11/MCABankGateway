package http

type paymentRequest struct {
	CustomerId     uint    `json:"customer_id" example:"1"`
	Amount         float64 `json:"amount" example:"100.50"`
	IdempotencyKey string  `json:"idempotency_key" example:"unique-key-123"`
}

type depositResponse struct {
	Status string `json:"status" example:"Depósito iniciado com sucesso"`
}

type withdrawResponse struct {
	Status string `json:"status" example:"Saque iniciado com sucesso"`
}
