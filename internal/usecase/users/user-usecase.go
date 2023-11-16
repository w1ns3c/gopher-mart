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
