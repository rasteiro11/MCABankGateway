package http

type createCustomerRequest struct {
	Nome  string `json:"nome" validate:"required" example:"John Doe"`
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

type updateCustomerRequest struct {
	Nome  string `json:"nome" validate:"required" example:"John Doe"`
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

type customerResponse struct {
	ID    uint    `json:"id" example:"1"`
	Nome  string  `json:"nome" example:"John Doe"`
	Email string  `json:"email" example:"john@example.com"`
	Saldo float64 `json:"saldo" example:"1000.50"`
}
