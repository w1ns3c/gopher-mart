package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"gopher-mart/internal/domain"
	"net/http"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

func CreateJWTcookie(userid uint64, secret string, lifetime time.Duration) (cookie *http.Cookie, err error) {
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
		Name:    domain.CookieName,
		Value:   jwtToken,
		Expires: cookieTime,

		HttpOnly: true,
		SameSite: 0,
	}
	return cookie, nil
}

func CheckJWTcookie(cookie *http.Cookie, secret string) (userID int64) {

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
		return domain.InvalidUserID
	}

	if !token.Valid {
		//fmt.Println("Token is not valid")
		return domain.InvalidUserID
	}

	// возвращаем ID пользователя в читаемом виде
	return int64(claims.UserID)

}
