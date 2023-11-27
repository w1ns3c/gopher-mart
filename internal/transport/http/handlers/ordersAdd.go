package handlers

import (
	"context"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/usecase/orders"
	usecaseUsers "gopher-mart/internal/usecase/users"
	"io"
	"net/http"
)

type OrdersAddHandler struct {
	usecase ordersAddUsecase
}

func NewOrdersAddHandler(usecase ordersAddUsecase) *OrdersAddHandler {
	return &OrdersAddHandler{usecase: usecase}
}

type ordersAddUsecase interface {
	//AddOrder(ctx context.Context, order *orders.Order) error
	AddOrder(ctx context.Context, user *users.User, orderNumber string) error
	orders.OrderValidator
	usecaseUsers.UserContextUsecase
}

func (h *OrdersAddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := h.usecase.CheckUserInContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("content-type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !h.usecase.ValidateOrderFormat(r.Context(), string(body)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.usecase.AddOrder(r.Context(), user, string(body))
	if err != nil {
		switch err {
		case errors.ErrAlreadyExist:
			w.WriteHeader(http.StatusOK)
		case errors.ErrCreatedByAnother:
			w.WriteHeader(http.StatusConflict)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
