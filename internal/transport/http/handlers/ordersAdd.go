package handlers

import (
	"context"
	"github.com/rs/zerolog/log"
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
		log.Err(err).Send()
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("content-type") != "text/plain" {
		log.Err(errors.ErrWrongContentType).Send()
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.usecase.AddOrder(r.Context(), user, string(body))
	if err != nil {
		log.Err(err).Send()
		switch err {
		case errors.ErrOrderAlreadyExist:
			w.WriteHeader(http.StatusOK)
		case errors.ErrOrderCreatedByAnother:
			w.WriteHeader(http.StatusConflict)
		case errors.ErrOrderWrongFormat:
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
