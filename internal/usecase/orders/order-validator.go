package orders

import (
	"context"
	"gopher-mart/internal/utils"
)

type OrderValidator interface {
	ValidateOrderFormat(ctx context.Context, orderNumber string) bool
	GetMaxRequestsPerMinute() uint64
	SetMaxRequestsPerMinute(max uint64)
}

type OrdersValidator struct {
	MaxRequests uint64
}

func NewOrdersValidator(max uint64) *OrdersValidator {
	return &OrdersValidator{MaxRequests: max}
}

func (validator *OrdersValidator) ValidateOrderFormat(ctx context.Context, orderNumber string) bool {
	return utils.LuhnValidator(orderNumber)
}

func (validator *OrdersValidator) GetMaxRequestsPerMinute() uint64 {
	return validator.MaxRequests
}

func (validator *OrdersValidator) SetMaxRequestsPerMinute(max uint64) {
	validator.MaxRequests = max
}
