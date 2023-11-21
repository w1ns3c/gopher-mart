package handlers

import (
	"context"
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
	GetBalance(ctx context.Context)
	usecaseUsers.UserContextUsecase
}

func (h *BalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := h.usecase.CheckUserInContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

}
