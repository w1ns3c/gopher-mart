package handlers

import (
	"context"
	"encoding/json"
	"gopher-mart/internal/domain/users"
	usecaseUsers "gopher-mart/internal/usecase/users"
	"net/http"
)

type BalanceHandler struct {
	usecase balanceUsecase
}

func NewBalanceHandler(usecase balanceUsecase) *BalanceHandler {
	return &BalanceHandler{usecase: usecase}
}

type balanceUsecase interface {
	GetBalance(ctx context.Context, user *users.User) (curBalance, withDrawn int, err error)
	usecaseUsers.UserContextUsecase
}
type responseBalance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

func (h *BalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := h.usecase.CheckUserInContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curBalance, withDrawn, err := h.usecase.GetBalance(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &responseBalance{
		Current:   curBalance,
		Withdrawn: withDrawn,
	}

	data, err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)

}
