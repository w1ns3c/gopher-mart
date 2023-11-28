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

type MarketUsecaseInf interface {
	usersUsecase.UserUsecase
	usersUsecase.UserBalanceUsecase
	usersUsecase.UserContextUsecase

	ordersUsecase.OrdersUsecaseInf
	ordersUsecase.OrderValidator
	cookies.CookiesUsecae
}

type GopherMart struct {
	cookies         cookies.Usecase
	users           usersUsecase.Usecase
	orders          ordersUsecase.Usecase
	ordersValidator ordersUsecase.OrdersValidator
}

func (g *GopherMart) LoginUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error) {
	return g.users.LoginUser(ctx, user)
}

func (g *GopherMart) RegisterUser(ctx context.Context, user *users.User) error {
	return g.users.RegisterUser(ctx, user)
}

func (g *GopherMart) GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error) {
	return g.users.GetUserWithdrawals(ctx, user)
}

func (g *GopherMart) CheckBalance(ctx context.Context, user *users.User) (curBalance, withDrawn int, err error) {
	return g.users.CheckBalance(ctx, user)
}

func (g *GopherMart) CheckUserInContext(ctx context.Context) (user *users.User, err error) {
	return g.users.CheckUserInContext(ctx)
}

func (g *GopherMart) ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error) {
	return g.orders.ListOrders(ctx, user)
}

func (g *GopherMart) AddOrder(ctx context.Context, user *users.User, orderNumber string) error {
	return g.orders.AddOrder(ctx, user, orderNumber)
}

func (g *GopherMart) CheckOrderStatus(ctx context.Context, orderNumber string) (order *orders.Order, err error) {
	return g.orders.CheckOrderStatus(ctx, orderNumber)
}

func (g *GopherMart) WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error {
	return g.orders.WithdrawBonuses(ctx, user, withdraw)
}

func (g *GopherMart) ValidateOrderFormat(ctx context.Context, orderNumber string) bool {
	return g.ordersValidator.ValidateOrderFormat(ctx, orderNumber)
}

func (g *GopherMart) GetMaxRequestsPerMinute() uint64 {
	return g.ordersValidator.GetMaxRequestsPerMinute()
}

func (g *GopherMart) SetMaxRequestsPerMinute(max uint64) {
	g.ordersValidator.SetMaxRequestsPerMinute(max)
}

func (g *GopherMart) ValidateCookie(ctx context.Context, cookie *http.Cookie) (user *users.User, err error) {
	return g.cookies.ValidateCookie(ctx, cookie)
}

func NewGophermart() *GopherMart {
	return &GopherMart{}
}
