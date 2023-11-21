package users

import (
	"context"
	"gopher-mart/internal/domain/users"
)

type UserBalanceUsecase interface {
	GetBalance(ctx context.Context, user *users.User) (curBalance, withDrawn int, err error)
}
