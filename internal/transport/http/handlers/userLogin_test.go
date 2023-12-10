package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type loginMock struct {
	Secret string
	Users  []loginRequest
}

func (l *loginMock) LoginUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error) {
	for _, mockUser := range l.Users {
		if user.Password == mockUser.Password &&
			user.Login == mockUser.Login {
			user.GenerateID(l.Secret)
			return utils.CreateJWTcookie(user.ID, l.Secret, time.Hour, domain.CookieName)
		}
	}
	return nil, fmt.Errorf("wrong credentials")
}

func TestLoginHandler_ServeHTTP(t *testing.T) {

	mock := &loginMock{
		Secret: "supersecret",
		Users: []loginRequest{
			{
				Login:    "user1",
				Password: "password",
			},
			{
				Login:    "user2",
				Password: "superpassword",
			},
		},
	}

	type args struct {
		usecase       loginUsecase
		credentials   *loginRequest
		headerContent string
	}
	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		// TODO: Add test cases.
		{
			name: "Test Valid user",
			args: args{usecase: mock,
				credentials: &loginRequest{
					Login:    "user1",
					Password: "password",
				},
				headerContent: "application/json",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Test Valid user2",
			args: args{usecase: mock,
				credentials: &loginRequest{
					Login:    "user2",
					Password: "superpassword",
				},
				headerContent: "application/json",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Test Wrong username",
			args: args{usecase: mock,
				credentials: &loginRequest{
					Login:    "user1111",
					Password: "password",
				},
				headerContent: "application/json",
			},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "Test Wrong password",
			args: args{usecase: mock,
				credentials: &loginRequest{
					Login:    "user2",
					Password: "password",
				},
				headerContent: "application/json",
			},
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "Test no content-type",
			args: args{usecase: mock,
				credentials: &loginRequest{
					Login:    "user2",
					Password: "password",
				},
				headerContent: "",
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name: "Test wrong content-type",
			args: args{usecase: mock,
				credentials: &loginRequest{
					Login:    "user2",
					Password: "password",
				},
				headerContent: "text/plain",
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &LoginHandler{
				usecase: tt.args.usecase,
			}
			body, err := json.Marshal(tt.args.credentials)
			if err != nil {
				t.Errorf("Can't mashal request body: %v", err)
				return
			}
			buf := bytes.NewBuffer(body)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8000", buf)
			req.Header.Add("content-type", tt.args.headerContent)
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)
			result := resp.Result()
			defer result.Body.Close()

			if req.Method != http.MethodPost {
				t.Errorf("LoginHandler got invalid HTTP method = %v , want = %v",
					req.Method, http.MethodPost)
				return
			}

			if result.StatusCode != tt.statusCode {
				t.Errorf("LoginHandler got status code = %v , want = %v",
					result.StatusCode, tt.statusCode)
				return
			}
			if tt.statusCode == http.StatusOK &&
				result.Header.Get("set-cookie") == "" {
				t.Errorf("LoginHandler don't add \"set-cookie\" header")
				return
			}

		})
	}
}
