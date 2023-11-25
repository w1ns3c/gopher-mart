package orders

import (
	"context"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
)

type OrdersUsecase interface {
	SendOrder()
	GetAllOrders()
	GetConfirmedOrders()
	GetDoneOrders()
	ConfirmOrder()
	DoneOrder()

	ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error)
	ValidateOrderFormat(ctx context.Context, orderNumber string) error
	AddOrder(ctx context.Context, user *users.User, orderNumber string) error
	WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error
}

type Usecase struct {
}

func (u *Usecase) ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error) {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) ValidateOrderFormat(ctx context.Context, orderNumber string) error {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) AddOrder(ctx context.Context, user *users.User, orderNumber string) error {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) WithdrawBonuses(ctx context.Context, user *users.User, orderNumber string, withdrawBonuses int) error {
	//TODO implement me
	panic("implement me")
}
