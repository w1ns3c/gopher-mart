package orders

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	locerrors "gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	"gopher-mart/internal/repository"
	"gopher-mart/internal/utils"
)

type OrdersUsecaseInf interface {
	ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error)
	AddOrder(ctx context.Context, user *users.User, orderNumber string) error
	WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error
}
type OrderValidator interface {
	ValidateOrderFormat(orderNumber string) bool
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

func (u *Usecase) ValidateOrderFormat(orderNumber string) bool {
	return utils.LuhnValidator(orderNumber)
}

func (u *Usecase) ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error) {
	return u.repo.ListOrders(ctx, user)
}
func (u *Usecase) AddOrder(ctx context.Context, user *users.User, orderNumber string) error {
	if !u.ValidateOrderFormat(orderNumber) {
		log.Error().Str("error", "wrong order number format")
		return locerrors.ErrOrderWrongFormat
	}
	_, userid, err := u.repo.CheckOrder(ctx, orderNumber)
	if err != nil {
		if errors.Is(err, locerrors.ErrWrongResultValues) {
			err = u.repo.AddOrder(ctx, user, orderNumber)
		}
		return err
	}

	// if not this user
	if userid != user.ID {
		return locerrors.ErrOrderCreatedByAnother
	}

	return locerrors.ErrOrderAlreadyExist
}

func (u *Usecase) WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error {
	err := u.repo.WithdrawBonuses(ctx, user, withdraw)
	/*
		1. Request user balance
		2. Compare balance with Sum (withdraw request)
		3. Sub: user_new_balance = user_balance - withdraw
		4. Save to DB
		4.1. Save: new balance for user
		4.2. Save:


	*/

	return err
}
