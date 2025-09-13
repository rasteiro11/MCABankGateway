package http

import "github.com/gofiber/fiber/v2"

type Handler interface {
	Deposit(c *fiber.Ctx) error
	Withdraw(c *fiber.Ctx) error
}
