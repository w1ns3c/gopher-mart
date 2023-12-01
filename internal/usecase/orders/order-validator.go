package orders

import (
	"gopher-mart/internal/utils"
)

type OrderValidator interface {
	ValidateOrderFormat(orderNumber string) bool
}

type OrdersValidator struct {
	MaxRequests uint64
}

func NewOrdersValidator() *OrdersValidator {
	return &OrdersValidator{}
}

func (validator *OrdersValidator) ValidateOrderFormat(orderNumber string) bool {
	return utils.LuhnValidator(orderNumber)
}
