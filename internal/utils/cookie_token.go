package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/domain/errors"
	"net/http"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func CreateJWTcookie(userid string, secret string, lifetime time.Duration, cookieName string) (cookie *http.Cookie, err error) {
	cookieTime := time.Now().Add(lifetime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(cookieTime),
		},
		UserID: userid,
	})
	jwtToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	cookie = &http.Cookie{
		Name:    cookieName,
		Value:   jwtToken,
		Expires: cookieTime,

		HttpOnly: true,
		SameSite: 0,
	}
	return cookie, nil
}

func CheckJWTcookie(cookie *http.Cookie, secret string) (userID string, err error) {

	// создаём экземпляр структуры с утверждениями
	claims := &Claims{}
	// парсим из строки токена tokenString в структуру claims
	token, err := jwt.ParseWithClaims(cookie.Value, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		log.Info().Msg("Token is not valid")
		return "", errors.ErrInvalidCookie
	}

	// возвращаем ID пользователя в читаемом виде
	return claims.UserID, nil

}
