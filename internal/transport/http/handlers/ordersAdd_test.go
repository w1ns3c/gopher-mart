package handlers

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ordersaddMock struct {
	orders []struct {
		userID  string
		orderID string
	}
}

func (o ordersaddMock) AddOrder(ctx context.Context, user *users.User, orderNumber string) error {
	// validate order
	// http.Status 422
	if !o.ValidateOrderFormat(orderNumber) {
		return errors.ErrOrderWrongFormat
	}

	for _, order := range o.orders {
		if order.orderID == orderNumber {
			if order.userID == user.ID {
				// http.StatusOK 200
				return errors.ErrOrderAlreadyExist
			}
			// http.StatusConflict 409
			return errors.ErrOrderCreatedByAnother
		}
	}
	return nil

}

func (o ordersaddMock) ValidateOrderFormat(orderNumber string) bool {
	return utils.LuhnValidator(orderNumber)
}

func (o ordersaddMock) CheckUserInContext(ctx context.Context) (user *users.User, err error) {
	anyType := ctx.Value(UserKey(domain.UserContextKey))
	userID, ok := anyType.(UserKey)
	if !ok {
		return nil, errors.ErrUserNotFoundInContext
	}
	user = new(users.User)
	user.ID = string(userID)
	return user, nil
}

func TestOrdersAddHandler_ServeHTTP(t *testing.T) {
	orders := []struct {
		userID  string
		orderID string
	}{
		{
			userID:  "userID1",
			orderID: "18",
		},
		{
			userID:  "userID2",
			orderID: "26",
		},
		{
			userID:  "userID3",
			orderID: "34",
		},
	}
	mock := &ordersaddMock{orders: orders}

	type args struct {
		usecase    ordersAddUsecase
		userID     string
		orderID    string
		statusCode int
	}
	tests := []struct {
		name string
		args *args
	}{
		// TODO: Add test cases.
		{
			name: "Test 1 Exist orderID, Valid userID",
			args: &args{
				usecase:    mock,
				userID:     "userID1",
				orderID:    "18",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Test 2 Exist orderID, Other userID",
			args: &args{
				usecase:    mock,
				userID:     "userID2",
				orderID:    "34",
				statusCode: http.StatusConflict,
			},
		},
		{
			name: "Test 3 Not exist orderID",
			args: &args{
				usecase:    mock,
				userID:     "userID3",
				orderID:    "42",
				statusCode: http.StatusAccepted,
			},
		},
		{
			name: "Test 4 Invalid orderID",
			args: &args{
				usecase:    mock,
				userID:     "userID3",
				orderID:    "4211",
				statusCode: http.StatusUnprocessableEntity,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &OrdersAddHandler{
				usecase: tt.args.usecase,
			}
			body := bytes.NewBufferString(tt.args.orderID)
			req := httptest.NewRequest(http.MethodPost,
				"http://localhost:8000/api/user/orders", body)
			ctx := context.WithValue(context.Background(),
				UserKey(domain.UserContextKey), UserKey(tt.args.userID))
			req = req.WithContext(ctx)
			req.Header.Set("content-type", "text/plain")
			resp := httptest.NewRecorder()

			h.ServeHTTP(resp, req)
			result := resp.Result()
			defer result.Body.Close()

			if req.Header.Get("content-type") != "text/plain" {
				assert.Equal(t, http.StatusBadRequest, result.StatusCode)
			}

			require.Equal(t, tt.args.statusCode, result.StatusCode)
		})
	}
}
