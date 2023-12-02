package users

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Login    string // nickname
	Password string
	Hash     string
	Cookie   string
}

func NewUser(login, password string) *User {
	return &User{
		Login:    login,
		Password: password,
	}
}

func (u *User) GenerateHash(salt string) error {
	pass := []byte(fmt.Sprintf("%s.%s.%s%s", salt, u.Password, salt, salt))
	hashBytes, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Hash = string(hashBytes)
	return err
}

func (u *User) CheckPasswordHash(salt string) bool {
	hashOfPass := []byte(fmt.Sprintf("%s.%s.%s%s", salt, u.Password, salt, salt))
	err := bcrypt.CompareHashAndPassword([]byte(u.Hash), hashOfPass)
	return err == nil
}

func (u *User) GenerateID(salt string) {
	hash := md5.Sum([]byte(fmt.Sprintf("%s.%s.%s", salt, u.Login, salt)))
	u.ID = hex.EncodeToString(hash[:])
}
