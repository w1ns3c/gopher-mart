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
}

type Usecase struct {
}

func (u *Usecase) ListOrders(ctx context.Context) (orders []orders.Order, err error) {
	//TODO implement me
	panic("implement me")
}