package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type currentBalanceMock struct {
	all map[UserKey]*users.Balance
}

func (mock *currentBalanceMock) CheckBalance(ctx context.Context, user *users.User) (balance *users.Balance, err error) {
	balance, ok := mock.all[UserKey(user.ID)]
	if !ok {
		return nil, fmt.Errorf("mock, user not exist")
	}
	return balance, nil
}

func (mock *currentBalanceMock) CheckUserInContext(ctx context.Context) (user *users.User, err error) {
	anyType := ctx.Value(UserKey(domain.UserContextKey))
	userID, ok := anyType.(UserKey)
	if !ok {
		return nil, errors.ErrUserNotFoundInContext
	}
	user = new(users.User)
	user.ID = string(userID)
	return user, nil
}

type UserKey string

func TestBalanceStatusHandler_ServeHTTP(t *testing.T) {
	users := map[UserKey]*users.Balance{
		"userID1": {Current: 100, WithdrawsSum: 0},
		"userID2": {Current: 100.222, WithdrawsSum: 11.356},
		"userID3": {Current: 100, WithdrawsSum: 300},
	}

	mock := &currentBalanceMock{all: users}

	type args struct {
		usecase    balanceUsecase
		userID     string
		balance    responseBalance
		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Valid user",
			args: args{
				usecase: mock,
				userID:  "userID2",
				balance: responseBalance{
					Current:   100.222,
					Withdrawn: 11.356,
				},
				statusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Test Wrong user",
			args: args{
				usecase:    mock,
				userID:     "userID2222",
				balance:    responseBalance{},
				statusCode: http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BalanceStatusHandler{
				usecase: tt.args.usecase,
			}

			req := httptest.NewRequest(http.MethodGet, "http://localhost:8000", nil)
			ctx := context.WithValue(context.Background(),
				UserKey(domain.UserContextKey), UserKey(tt.args.userID))
			req = req.WithContext(ctx)
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)
			result := resp.Result()
			defer result.Body.Close()

			require.Equal(t, tt.args.statusCode, result.StatusCode, "compare status codes")
			var balance responseBalance
			body, err := io.ReadAll(result.Body)
			if err != nil && !tt.wantErr {
				t.Errorf("got error = %v", err)
				return
			}

			if !tt.wantErr {
				require.NoError(t, err, "Can't get json from raw body")
				err = json.Unmarshal(body, &balance)
				require.NoError(t, err, "Can't get balance from json body")

				require.Equal(t, tt.args.balance.Current, balance.Current,
					"current user balance")

				require.Equal(t, tt.args.balance.Withdrawn, balance.Withdrawn,
					"current user balance")
				require.Equal(t, tt.args.statusCode, result.StatusCode, "compare status codes")
			}
		})
	}
}
