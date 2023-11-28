package repository

import (
	"context"
	"database/sql"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	"net/http"
)

type UserRepoInf interface {
	Init(ctx context.Context) error
	LoginUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error)
	RegisterUser(ctx context.Context, user *users.User) error
	GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error)
	CheckBalance(ctx context.Context, user *users.User) (curBalance, withDrawn int, err error)

	GetUser(ctx context.Context, id string) (user *users.User, err error)
}

type UserRepo struct {
	db  *sql.DB
	url string
}
