package postgres

import (
	"context"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	"net/http"
)

func (pg *PostgresRepo) ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error) {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresRepo) AddOrder(ctx context.Context, user *users.User, orderNumber string) error {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresRepo) WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresRepo) CheckOrderStatus(ctx context.Context, orderNumber string) (order *orders.Order, err error) {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresRepo) GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error) {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresRepo) CheckBalance(ctx context.Context, user *users.User) (curBalance, withDrawn int, err error) {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresRepo) CheckUserExist(ctx context.Context, cookie *http.Cookie) (user *users.User, err error) {
	//TODO implement me
	panic("implement me")
}
