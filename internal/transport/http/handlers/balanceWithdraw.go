package handlers

import (
	"context"
	"encoding/json"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	usecaseUsers "gopher-mart/internal/usecase/users"
	"gopher-mart/internal/utils"
	"net/http"
)

type balanceWithdrawHandler struct {
	usecase balanceWithdrawUsecase
}

func NewBalanceWithdrawHandler(usecase balanceWithdrawUsecase) *balanceWithdrawHandler {
	return &balanceWithdrawHandler{usecase: usecase}
}

type balanceWithdrawUsecase interface {
	WithdrawBonuses(ctx context.Context, user *users.User, orderNumber string, withdrawBonuses int) error
	usecaseUsers.UserContextUsecase
}

type withdrawRequest struct {
	OrderID string `json:"order"`
	Sum     int    `json:"sum"`
}

func (h *balanceWithdrawHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := h.usecase.CheckUserInContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req *withdrawRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !utils.LuhnValidator(req.OrderID) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = h.usecase.WithdrawBonuses(r.Context(), user, req.OrderID, req.Sum)
	if err != nil {
		switch err {
		case errors.ErrNotEnoughBonuses:
			w.WriteHeader(http.StatusPaymentRequired)
		case errors.ErrWrongOrder:
			w.WriteHeader(http.StatusUnprocessableEntity)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
