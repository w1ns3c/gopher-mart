package users

import (
	"context"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	"gopher-mart/internal/repository"
	"gopher-mart/internal/utils"
	"net/http"
	"time"
)

type UserUsecase interface {
	LoginUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error)
	RegisterUser(ctx context.Context, user *users.User) error

	GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error)
}

type Usecase struct {
	storage        repository.UsersRepoInf
	Secret         string
	CookieName     string
	CookieLifeTime time.Duration
}

func NewUsecase(storage repository.UsersRepoInf, secret,
	cookieName string, lifetime time.Duration) *Usecase {
	return &Usecase{
		Secret:         secret,
		CookieName:     cookieName,
		CookieLifeTime: lifetime,
		storage:        storage,
	}
}

func (u *Usecase) CheckUserInContext(ctx context.Context) (user *users.User, err error) {
	anyType := ctx.Value(domain.UserContextKey)
	user, ok := anyType.(*users.User)
	if !ok {
		return nil, errors.ErrUserNotFoundInContext
	}
	return user, nil
}

func (u *Usecase) CheckBalance(ctx context.Context,
	user *users.User) (curBalance, withDrawn int, err error) {
	return u.storage.CheckBalance(ctx, user)
}

func (u *Usecase) LoginUser(ctx context.Context,
	user *users.User) (cookie *http.Cookie, err error) {
	hash, err := u.storage.LoginUser(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Hash = hash
	if !user.CheckPasswordHash(u.Secret) {
		return nil, errors.ErrUserLogin
	}
	return utils.CreateJWTcookie(user.ID, u.Secret, u.CookieLifeTime, u.CookieName)
}

func (u *Usecase) RegisterUser(ctx context.Context, user *users.User) error {
	err := user.GenerateHash(u.Secret)
	if err != nil {
		return err
	}
	return u.storage.RegisterUser(ctx, user)
}

func (u *Usecase) GetUserWithdrawals(ctx context.Context,
	user *users.User) (wd []withdraws.Withdraw, err error) {
	return u.storage.GetUserWithdrawals(ctx, user)
}
