package repository

import (
	"context"
	"database/sql"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	"net/http"
)

type UsersRepoInf interface {
	LoginUser(ctx context.Context, user *users.User) (userID uint64, userHash string, err error)
	RegisterUser(ctx context.Context, user *users.User) error
	GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error)
	CheckBalance(ctx context.Context, user *users.User) (curBalance, withDrawn int, err error)

	CheckUserExist(ctx context.Context, cookie *http.Cookie) (user *users.User, err error)
}

type UserRepo struct {
	db  *sql.DB
	url string
}
