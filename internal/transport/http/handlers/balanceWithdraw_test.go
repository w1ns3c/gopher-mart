package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	"gopher-mart/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

type withdrawBalanceMock struct {
	all map[UserKey]*users.Balance
}

func (m *withdrawBalanceMock) WithdrawUserBonuses(ctx context.Context,
	user *users.User, wd *withdraws.Withdraw) error {
	balance, ok := m.all[UserKey(user.ID)]
	if !ok {
		return fmt.Errorf("user not found in mock")
	}

	if balance.Current < wd.Sum {
		return errors.ErrNotEnoughBonuses
	}

	if !utils.LuhnValidator(wd.OrderID) {
		return errors.ErrOrderWrongFormat
	}
	return nil
}

func (m *withdrawBalanceMock) CheckUserInContext(ctx context.Context) (user *users.User, err error) {
	anyType := ctx.Value(UserKey(domain.UserContextKey))
	userID, ok := anyType.(UserKey)
	if !ok {
		return nil, errors.ErrUserNotFoundInContext
	}
	user = new(users.User)
	user.ID = string(userID)
	return user, nil
}

func Test_balanceWithdrawHandler_ServeHTTP(t *testing.T) {
	users := map[UserKey]*users.Balance{
		"userID1": {Current: 200.54, WithdrawsSum: 0},
		"userID2": {Current: 100.222, WithdrawsSum: 11.356},
		"userID3": {Current: 100, WithdrawsSum: 300},
	}

	mock := &withdrawBalanceMock{all: users}

	type args struct {
		usecase    balanceWithdrawUsecase
		userID     string
		withdraw   withdrawRequest
		statusCode int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Valid user and orderID, correct withdraw sum",
			args: args{
				usecase: mock,
				userID:  "userID2",
				withdraw: withdrawRequest{
					Sum:     100.222,
					OrderID: "52284247",
				},
				statusCode: http.StatusOK,
			},
			wantErr: false,
		},
		{
			name: "Test Valid user and orderID, withdraw sum TOO big",
			args: args{
				usecase: mock,
				userID:  "userID2",
				withdraw: withdrawRequest{
					Sum:     200.55,
					OrderID: "52284247",
				},
				statusCode: http.StatusPaymentRequired,
			},
			wantErr: true,
		},
		{
			name: "Test Invalid orderID",
			args: args{
				usecase: mock,
				userID:  "userID2",
				withdraw: withdrawRequest{
					Sum:     100.55,
					OrderID: "5228424",
				},
				statusCode: http.StatusPaymentRequired,
			},
			wantErr: true,
		},
		{
			name: "Test Invalid userID",
			args: args{
				usecase: mock,
				userID:  "userID2222",
				withdraw: withdrawRequest{
					Sum:     100.55,
					OrderID: "5228424",
				},
				statusCode: http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &balanceWithdrawHandler{
				usecase: tt.args.usecase,
			}

			body, err := json.Marshal(tt.args.withdraw)
			require.NoError(t, err, "Can't marshal withdraw to json body")
			buf := bytes.NewBuffer(body)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8000", buf)
			ctx := context.WithValue(context.Background(),
				UserKey(domain.UserContextKey), UserKey(tt.args.userID))
			req = req.WithContext(ctx)
			req.Header.Set("content-type", "application/json")
			resp := httptest.NewRecorder()
			h.ServeHTTP(resp, req)
			result := resp.Result()
			defer result.Body.Close()

			require.Equal(t, "application/json", req.Header.Get("content-type"))
			require.Equal(t, tt.args.statusCode, result.StatusCode, "compare status codes")

		})
	}
}
