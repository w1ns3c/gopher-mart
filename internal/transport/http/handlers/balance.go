package handlers

import "context"

type BalanceHandler struct {
	usecase balanceUsecase
}

func NewBalanceHandler(usecase balanceUsecase) *BalanceHandler {
	return &BalanceHandler{usecase: usecase}
}

type balanceUsecase interface {
	GetBalance(ctx context.Context)
}
