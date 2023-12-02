package repository

import (
	"context"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
)

type Repository interface {
	Init(ctx context.Context) error
	CheckConnection() error
	OrdersRepoInf
	UsersRepoInf
}

type OrdersRepoInf interface {
	ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error)
	AddOrder(ctx context.Context, user *users.User, orderNumber string) error
	WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error
	CheckOrder(ctx context.Context, orderNumber string) (orderid, userid string, err error)
}

type UsersRepoInf interface {
	LoginUser(ctx context.Context, user *users.User) (userID string, userHash string, err error)
	RegisterUser(ctx context.Context, user *users.User) error
	GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error)
	CheckUserExist(ctx context.Context, userID string) (user *users.User, err error)

	BalanceRepoInf
}

type BalanceRepoInf interface {
	CheckBalance(ctx context.Context, user *users.User) (curBalance, withDrawn uint64, err error)
	UpdateBalance(ctx context.Context, user *users.User, balance *users.Balance) error
}
