package repository

import (
	"context"
	"gopher-mart/internal/domain/users"
)

type userRepo interface {
	GetUser(ctx context.Context, id string) (user *users.User, err error)
	SaveUser(ctx context.Context, user *users.User) error
}
