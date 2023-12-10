package middlewares

import (
	"context"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

type authMock struct {
	Secret string
}

func (mock *authMock) ValidateCookie(ctx context.Context, cookie *http.Cookie) (user *users.User, err error) {
	userID, err := utils.CheckJWTcookie(cookie, mock.Secret)
	return &users.User{ID: userID}, err
}

func TestAuthMiddleware_AuthMiddleware(t *testing.T) {

	authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	mock := &authMock{Secret: "supersecret"}

	type args struct {
		usecase AuthUsecase
		cookie  *http.Cookie
	}

	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		// TODO: Add test cases.
		{
			name: "Test valid cookie",
			args: args{
				usecase: mock,
				cookie: &http.Cookie{
					Name:  domain.CookieName,
					Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE0ODk1ODYxMTEsIlVzZXJJRCI6IjYzZjFjMzlmODY3YjI5NDA0NGQ0OGRjZTYwYzdkY2JjIn0.5_ibfJ6AwWnQO0qWzxYoDXLwfGKrLx59aLI48MCvhDU",
				},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Test expired cookie",
			args: args{
				usecase: mock,
				cookie: &http.Cookie{
					Name:  domain.CookieName,
					Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDE4NjExMSwiVXNlcklEIjoiMTIzMjEzMjExIn0.NGAPdTMreBv5V_mR7bEq1vEMsecpgQBopVEParQBe1g",
				},
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "Test no cookie",
			args: args{
				usecase: mock,
				cookie:  &http.Cookie{},
			},
			wantCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &AuthMiddleware{
				usecase: tt.args.usecase,
			}
			req := httptest.NewRequest(http.MethodGet,
				"http://localhost:8000", nil)

			req.AddCookie(tt.args.cookie)
			res := httptest.NewRecorder()
			authHandler(res, req)

			got := m.AuthMiddleware(authHandler)

			got.ServeHTTP(res, req)
			result := res.Result()
			defer result.Body.Close()

			if result.StatusCode != tt.wantCode {
				t.Errorf("AuthMiddleware() wrong status code = %v, want %v",
					result.StatusCode, tt.wantCode)
				return
			}
		})
	}
}
