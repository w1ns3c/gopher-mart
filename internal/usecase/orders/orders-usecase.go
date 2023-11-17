package orders

import (
	"context"
	"gopher-mart/internal/domain/orders"
)

type OrdersUsecase interface {
	SendOrder()
	GetAllOrders()
	GetConfirmedOrders()
	GetDoneOrders()
	ConfirmOrder()
	DoneOrder()

	ListOrders(ctx context.Context) (orders []orders.Order, err error)
	ValidateOrderFormat(ctx context.Context, orderNumber string) error
	AddOrder(ctx context.Context, orderNumber string) error
}

type Usecase struct {
}

func (u *Usecase) ListOrders(ctx context.Context) (orders []orders.Order, err error) {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) ValidateOrderFormat(ctx context.Context, orderNumber string) error {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) AddOrder(ctx context.Context, orderNumber string) error {
	//TODO implement me
	panic("implement me")
}
