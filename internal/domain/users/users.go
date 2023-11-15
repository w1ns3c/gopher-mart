package users

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type ID string

type User struct {
	ID       ID     `json:"id"`
	Password string `json:"password"`
	Hash     string `json:"-"`
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
