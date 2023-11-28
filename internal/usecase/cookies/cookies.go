package cookies

import (
	"context"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/repository"
	"net/http"
)

type CookiesUsecae interface {
	ValidateCookie(ctx context.Context, cookie *http.Cookie) (user *users.User, err error)
}

type Usecase struct {
	storage repository.UserRepoInf
}

func (u Usecase) ValidateCookie(ctx context.Context, cookie *http.Cookie) (user *users.User, err error) {
	return u.storage.ValidateCookie(ctx, cookie)
}
