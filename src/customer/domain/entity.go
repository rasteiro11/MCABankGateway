package domain

type Customer struct {
	ID    uint    `json:"id"`
	Nome  string  `json:"nome"`
	Email string  `json:"email"`
	Saldo float64 `json:"saldo"`
}
