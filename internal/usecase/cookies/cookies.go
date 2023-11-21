package cookies

import (
	"context"
	"gopher-mart/internal/domain/users"
	"net/http"
)

type CookiesUsecae interface {
	ValidateCookie(ctx context.Context, cookie *http.Cookie) (user *users.User, err error)
}

type Usecase struct {
}

func (u Usecase) ValidateCookie(ctx context.Context, cookie *http.Cookie) (user *users.User, err error) {
	//TODO implement me
	panic("implement me")
}
