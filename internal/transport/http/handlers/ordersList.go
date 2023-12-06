package handlers

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/domain/users"
	usecaseUsers "gopher-mart/internal/usecase/users"
	"net/http"
	"time"
)

type ordersListHandler struct {
	usecase ordersListUsecase
}

func NewOrdersListHandler(usecase ordersListUsecase) *ordersListHandler {
	return &ordersListHandler{usecase: usecase}
}

type ordersListUsecase interface {
	// TODO orders MUST be sorted by date, should it be on repo level?
	ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error)
	usecaseUsers.UserContextUsecase
}

type ordersResponse struct {
	ID      string             `json:"number"`
	Status  orders.OrderStatus `json:"status,omitempty"`
	Accrual float64            `json:"accrual,omitempty"` // accrual
	Date    time.Time          `json:"-"`
}

func (r *ordersResponse) MarshalJSON() ([]byte, error) {
	type Alias ordersResponse
	if r.Date.IsZero() {
		return json.Marshal(&struct {
			*Alias
		}{
			(*Alias)(r),
		})
	}
	return json.Marshal(&struct {
		Date string `json:"uploaded_at,omitempty"`
		*Alias
	}{
		r.Date.Format(time.RFC3339),
		(*Alias)(r),
	})
}

func (h *ordersListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	user, err := h.usecase.CheckUserInContext(r.Context())
	if err != nil {
		log.Err(err).Send()
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodGet {
		log.Err(errors.ErrMethodNotAllowed).Send()
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	orders, err := h.usecase.ListOrders(r.Context(), user)
	if err != nil {
		log.Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// orders not found
	l := len(orders)
	if l == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	resp := make([]ordersResponse, l)
	for id, order := range orders {
		resp[id].ID = order.ID
		resp[id].Status = order.Status
		resp[id].Accrual = order.Cashback
		resp[id].Date = order.Date
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
