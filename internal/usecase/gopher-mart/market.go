package gopher_mart

import (
	"gopher-mart/internal/usecase/orders"
	"gopher-mart/internal/usecase/users"
)

type MarketUsecase interface {
	users.UserUsecase
	orders.OrdersUsecase
}
