package http

import (
	"github.com/go-chi/chi/v5"
	"gopher-mart/internal/transport/http/handlers"
	gopher_mart "gopher-mart/internal/usecase/gopher-mart"
)

func NewRouter(market gopher_mart.MarketUsecase) error {
	// init handlers
	loginHandler := handlers.NewLoginHandler(market)
	registerHandler := handlers.NewRegisterHandler(market)

	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", loginHandler.ServeHTTP)
			r.Post("/login", registerHandler.ServeHTTP)

		})
	})
}
