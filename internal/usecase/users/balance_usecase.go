package users

import (
	"context"
	"gopher-mart/internal/domain/users"
)

type UserBalanceUsecase interface {
	CheckBalance(ctx context.Context, user *users.User) (balance *users.Balance, err error)
	UpdateBalance(ctx context.Context, user *users.User, balance *users.Balance) error
}

type UserContextUsecase interface {
	CheckUserInContext(ctx context.Context) (user *users.User, err error)
}

func (u *Usecase) CheckBalance(ctx context.Context, user *users.User) (balance *users.Balance, err error) {
	current, withdraws, err := u.repo.CheckBalance(ctx, user)
	balance = &users.Balance{
		Current:      current,
		WithdrawsSum: withdraws,
	}
	return balance, err
}

func (u *Usecase) UpdateBalance(ctx context.Context, user *users.User, balance *users.Balance) error {
	return u.repo.UpdateBalance(ctx, user, balance)
}
