package repository

import (
	"context"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
)

type OrdersRepoInf interface {
	ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error)
	AddOrder(ctx context.Context, user *users.User, orderNumber string) error
	WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error
	CheckOrderStatus(ctx context.Context, orderNumber string) (order *orders.Order, err error)
}
