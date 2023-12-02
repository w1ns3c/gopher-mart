package users

import (
	"context"
	"gopher-mart/internal/domain/users"
)

type UserBalanceUsecase interface {
	CheckBalance(ctx context.Context, user *users.User) (curBalance, withDrawn uint64, err error)
}

type UserContextUsecase interface {
	CheckUserInContext(ctx context.Context) (user *users.User, err error)
}
