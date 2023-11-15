package handlers

import (
	"encoding/json"
	dusers "gopher-mart/internal/domain/users"
	"gopher-mart/internal/usecase/users"
	"net/http"
)

type LoginHandler struct {
	loginUsecase users.UserUsecase
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("content-type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user dusers.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie, err := h.loginUsecase.LoginUser(user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("set-cookie", cookie)
	w.Write([]byte("OK"))
}
