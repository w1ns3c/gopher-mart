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

type Options func(usecase *Usecase)

func NewUsecaseWith(options ...Options) *Usecase {
	usecase := new(Usecase)
	for _, option := range options {
		option(usecase)
	}
	return usecase
}

func WithRepo(repo repository.OrdersRepoInf) func(u *Usecase) {
	return func(u *Usecase) {
		u.repo = repo
	}
}

type Usecase struct {
	repo repository.OrdersRepoInf
}

func (u *Usecase) ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error) {
	return u.repo.ListOrders(ctx, user)
}

func (u *Usecase) AddOrder(ctx context.Context, user *users.User, orderNumber string) error {
	return u.repo.AddOrder(ctx, user, orderNumber)
}

func (u *Usecase) WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error {
	return u.repo.WithdrawBonuses(ctx, user, withdraw)
}

func (u *Usecase) CheckOrderStatus(ctx context.Context, orderNumber string) (order *orders.Order, err error) {
	return u.repo.CheckOrderStatus(ctx, orderNumber)
}
