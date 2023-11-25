package users

import (
	"context"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
)

type UserUsecase interface {
	LoginUser(ctx context.Context, user *users.User) (cookie string, err error)
	RegisterUser(ctx context.Context, user *users.User) error

	GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error)

	UserContextUsecase
}
type Usecase struct {
}

func (u *Usecase) CheckUserInContext(ctx context.Context) (user *users.User, err error) {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) LoginUser(ctx context.Context, user *users.User) (cookie string, err error) {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) RegisterUser(ctx context.Context, user *users.User) error {
	//TODO implement me
	panic("implement me")
}

func NewUsecase() *Usecase {
	return &Usecase{}
}

type UserContextUsecase interface {
	CheckUserInContext(ctx context.Context) (user *users.User, err error)
}

func (u *Usecase) GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error) {
	//TODO implement me
	panic("implement me")
}
