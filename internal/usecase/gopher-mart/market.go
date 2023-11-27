package gophermart

import (
	"gopher-mart/internal/usecase/cookies"
	"gopher-mart/internal/usecase/orders"
	"gopher-mart/internal/usecase/users"
)

type MarketUsecase interface {
	users.UserUsecase
	users.UserBalanceUsecase
	orders.OrdersUsecase
	orders.OrderValidator
	cookies.CookiesUsecae
}

type Usecase struct {
}
