package gophermart

import (
	"context"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	"gopher-mart/internal/usecase/cookies"
	ordersUsecase "gopher-mart/internal/usecase/orders"
	usersUsecase "gopher-mart/internal/usecase/users"
	"net/http"
)

type MarketUsecase interface {
	usersUsecase.UserUsecase
	usersUsecase.UserBalanceUsecase
	usersUsecase.UserContextUsecase

	ordersUsecase.OrdersUsecase
	ordersUsecase.OrderValidator
	cookies.CookiesUsecae
}

type GopherMart struct {
	usersUsecase.Usecase
	ordersUsecase.OrdersValidator
}

func NewGophermart() *GopherMart {
	return &GopherMart{}
}

func (u GopherMart) ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error) {
	//TODO implement me
	panic("implement me")
}

func (u GopherMart) AddOrder(ctx context.Context, user *users.User, orderNumber string) error {
	//TODO implement me
	panic("implement me")
}

func (u GopherMart) CheckOrderStatus(ctx context.Context, orderNumber string) (order *orders.Order, err error) {
	//TODO implement me
	panic("implement me")
}

func (u GopherMart) WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error {
	//TODO implement me
	panic("implement me")
}

func (u GopherMart) ValidateOrderFormat(ctx context.Context, orderNumber string) bool {
	//TODO implement me
	panic("implement me")
}

func (u GopherMart) GetMaxRequestsPerMinute() uint64 {
	//TODO implement me
	panic("implement me")
}

func (u GopherMart) SetMaxRequestsPerMinute(max uint64) {
	//TODO implement me
	panic("implement me")
}

func (u GopherMart) ValidateCookie(ctx context.Context, cookie *http.Cookie) (user *users.User, err error) {
	//TODO implement me
	panic("implement me")
}
