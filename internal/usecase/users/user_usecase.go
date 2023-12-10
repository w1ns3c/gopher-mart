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

type UserUsecaseInf interface {
	LoginUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error)
	RegisterUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error)

	GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error)
	UserBalanceUsecase
	UserContextUsecase
}

type Usecase struct {
	repo           repository.UsersRepoInf
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
		repo:           storage,
	}
}

type Options func(usecase *Usecase)

func NewUsecaseWith(options ...Options) *Usecase {
	usecase := new(Usecase)
	for _, option := range options {
		option(usecase)
	}
	return usecase
}

func WithSecret(secret string) func(u *Usecase) {
	return func(u *Usecase) {
		u.Secret = secret
	}
}

func WithCookieName(cookieName string) func(u *Usecase) {
	return func(u *Usecase) {
		u.CookieName = cookieName
	}
}

func WithCookieLifetime(cookieLife time.Duration) func(u *Usecase) {
	return func(u *Usecase) {
		u.CookieLifeTime = cookieLife
	}
}

func WithRepo(repo repository.UsersRepoInf) func(u *Usecase) {
	return func(u *Usecase) {
		u.repo = repo
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

func (u *Usecase) LoginUser(ctx context.Context,
	user *users.User) (cookie *http.Cookie, err error) {
	userID, hash, err := u.repo.LoginUser(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Hash = hash
	if !user.CheckPasswordHash(u.Secret) {
		return nil, errors.ErrUserLogin
	}
	user.ID = userID
	return utils.CreateJWTcookie(user.ID, u.Secret, u.CookieLifeTime, u.CookieName)
}

func (u *Usecase) RegisterUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error) {
	err = user.GenerateHash(u.Secret)
	if err != nil {
		return nil, err
	}
	user.GenerateID(u.Secret)
	err = u.repo.RegisterUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return utils.CreateJWTcookie(user.ID, u.Secret, u.CookieLifeTime, u.CookieName)
}

func (u *Usecase) GetUserWithdrawals(ctx context.Context,
	user *users.User) (wd []withdraws.Withdraw, err error) {
	return u.repo.GetUserWithdrawals(ctx, user)
}
