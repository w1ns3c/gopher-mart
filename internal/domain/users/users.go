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
	Cookie   string
}

func (u *User) GenerateHash(salt string) error {
	pass := []byte(fmt.Sprintf("%s.%s.%s%s", salt, u.Password, salt, salt))
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Hash = string(hash)
	return nil
}

func (u *User) CheckPasswordHash(salt string) bool {
	pass := []byte(fmt.Sprintf("%s.%s.%s%s", salt, u.Password, salt, salt))
	err := bcrypt.CompareHashAndPassword([]byte(u.Hash), pass)
	return err == nil
}

func NewUser(login, password string) *User {
	return &User{
		Login:    login,
		Password: password,
	}
}
