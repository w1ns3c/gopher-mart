package users

import "gopher-mart/internal/domain/users"

type UserUsecase interface {
	LoginUser(users.User) (cookie string, err error)
	RegisterUser(users.User) error
}
