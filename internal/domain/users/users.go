package users

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint64
	Login    string // nickname
	Password string
	Hash     string
}

func (u *User) GenerateHash(salt string) (string, error) {
	pass := []byte(fmt.Sprintf("%s.%s.%s", salt, u.Password, salt))
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	return string(hash), err
}

func (u *User) CheckPasswordHash(hash, salt string) bool {
	pass := []byte(fmt.Sprintf("%s.%s.%s", salt, u.Password, salt))
	err := bcrypt.CompareHashAndPassword([]byte(hash), pass)
	return err == nil
}

func NewUser(login, password string) *User {
	return &User{
		Login:    login,
		Password: password,
	}
}
