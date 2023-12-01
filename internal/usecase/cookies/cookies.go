package cookies

import (
	"context"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/repository"
	"gopher-mart/internal/utils"
	"net/http"
)

type CookiesUsecae interface {
	ValidateCookie(ctx context.Context, cookie *http.Cookie) (user *users.User, err error)
}

type Usecase struct {
	Secret string
	repo   repository.UsersRepoInf
}

func (u *Usecase) ValidateCookie(ctx context.Context, cookie *http.Cookie) (user *users.User, err error) {
	userID := utils.CheckJWTcookie(cookie, u.Secret)
	if userID == domain.InvalidUserID {
		return nil, errors.ErrInvalidValue
	}
	user, err = u.repo.CheckUserExist(ctx, cookie)
	if err != nil {
		return nil, err
	}

	return user, nil

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

func WithRepo(repo repository.UsersRepoInf) func(u *Usecase) {
	return func(u *Usecase) {
		u.repo = repo
	}
}
