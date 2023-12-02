package cookies

import (
	"context"
	"github.com/rs/zerolog/log"
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
	userID, err := utils.CheckJWTcookie(cookie, u.Secret)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, err
	}

	user, err = u.repo.CheckUserExist(ctx, userID)
	if err != nil {
		log.Error().Err(err).Send()
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
