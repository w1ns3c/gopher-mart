package http

import (
	"github.com/go-chi/chi/v5"
	"gopher-mart/internal/transport/http/handlers"
	market "gopher-mart/internal/usecase/gopher-mart"
)

func NewRouter(market market.MarketUsecase) error {
	// init handlers
	loginHandler := handlers.NewLoginHandler(market)
	registerHandler := handlers.NewRegisterHandler(market)

	// orders
	listOrdersHandler := handlers.NewOrdersListHandler(market)

	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", loginHandler.ServeHTTP)
			r.Post("/register", registerHandler.ServeHTTP)

			r.Get("/orders", listOrdersHandler.ServeHTTP)

		})
	})
}
