package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type registerMock struct {
	Secret string
	Users  []reqisterRequest
}

func (m *registerMock) RegisterUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error) {
	var found = false
	for _, mockUser := range m.Users {
		if user.Login == mockUser.Login {
			found = true
			break
		}
	}

	// user already exist
	if found {
		return nil, fmt.Errorf("user already exist, duplicate key value")
	}

	// user not exist, create
	user.GenerateID(m.Secret)
	return utils.CreateJWTcookie(user.ID, m.Secret, time.Hour, domain.CookieName)
}

func TestRegisterHandler_ServeHTTP(t *testing.T) {
	mock := &registerMock{
		Secret: "supersecret",
		Users: []reqisterRequest{
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
		usecase     registerUsecase
		credentials reqisterRequest
		contentType string
	}
	tests := []struct {
		name       string
		args       args
		statusCode int
	}{
		// TODO: Add test cases.
		{
			name: "Test new user",
			args: args{
				usecase: mock,
				credentials: reqisterRequest{
					Login:    "newuser",
					Password: "newpassword",
				},
				contentType: "application/json",
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Test existing user",
			args: args{
				usecase: mock,
				credentials: reqisterRequest{
					Login:    "user1",
					Password: "newpassword",
				},
				contentType: "application/json",
			},
			statusCode: http.StatusConflict,
		},
		{
			name: "Test no content-type",
			args: args{
				usecase: mock,
				credentials: reqisterRequest{
					Login:    "user1",
					Password: "newpassword",
				},
				contentType: "",
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &RegisterHandler{
				usecase: tt.args.usecase,
			}

			body, err := json.Marshal(tt.args.credentials)
			if err != nil {
				t.Errorf("Can't mashal request body: %v", err)
				return
			}
			buf := bytes.NewBuffer(body)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8000", buf)
			req.Header.Add("content-type", tt.args.contentType)
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)
			result := resp.Result()
			defer result.Body.Close()

			if tt.statusCode != http.StatusBadRequest {
				assert.Equal(t, "application/json", req.Header.Get("content-type"),
					"RegisterHandler got invalid content-type in request header")
			}

			require.Equal(t, http.MethodPost, req.Method,
				"RegisterHandler got invalid HTTP method ")

			require.Equal(t, tt.statusCode, result.StatusCode,
				"RegisterHandler got invalid HTTP status code")

			if tt.statusCode == http.StatusOK {
				require.NotEmpty(t, result.Header.Get("set-cookie"),
					"RegisterHandler don't set cookie to response")
			}

		})
	}
}
