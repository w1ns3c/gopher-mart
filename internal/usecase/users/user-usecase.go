package users

import (
	"context"
	"gopher-mart/internal/domain/users"
)

type UserUsecase interface {
	LoginUser(ctx context.Context, user *users.User) (cookie string, err error)
	RegisterUser(ctx context.Context, user *users.User) error
}

type Usecase struct {
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
