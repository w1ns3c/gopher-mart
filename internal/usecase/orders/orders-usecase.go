package orders

import (
	"context"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	"gopher-mart/internal/repository"
)

type OrdersUsecaseInf interface {
	ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error)
	AddOrder(ctx context.Context, user *users.User, orderNumber string) error
	CheckOrderStatus(ctx context.Context, orderNumber string) (order *orders.Order, err error)
	WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error
}

type Usecase struct {
	storage repository.OrdersRepoInf
}

func (u *Usecase) ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error) {
	return u.storage.ListOrders(ctx, user)
}

func (u *Usecase) AddOrder(ctx context.Context, user *users.User, orderNumber string) error {
	return u.storage.AddOrder(ctx, user, orderNumber)
}

func (u *Usecase) WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error {
	return u.storage.WithdrawBonuses(ctx, user, withdraw)
}

func (u *Usecase) CheckOrderStatus(ctx context.Context, orderNumber string) (order *orders.Order, err error) {
	return u.storage.CheckOrderStatus(ctx, orderNumber)
}
