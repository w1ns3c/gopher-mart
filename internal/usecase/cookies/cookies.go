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
	Secret  string
	storage repository.UserRepoInf
}

func (u *Usecase) ValidateCookie(ctx context.Context, cookie *http.Cookie) (user *users.User, err error) {
	userID := utils.CheckJWTcookie(cookie, u.Secret)
	if userID == domain.InvalidUserID {
		return nil, errors.ErrInvalidValue
	}
	user, err = u.storage.CheckUserExist(ctx, cookie)
	if err != nil {
		return nil, err
	}

	return user, nil

}
