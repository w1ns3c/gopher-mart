package http

import (
	"github.com/go-chi/chi/v5"
	"gopher-mart/internal/transport/http/handlers"
	market "gopher-mart/internal/usecase/gopher-mart"
	"net/http"
)

func NewRouter(market market.MarketUsecase) http.Handler {
	// init handlers
	loginHandler := handlers.NewLoginHandler(market)
	registerHandler := handlers.NewRegisterHandler(market)

	// orders
	listOrdersHandler := handlers.NewOrdersListHandler(market)
	addOrdersHandler := handlers.NewOrdersAddHandler(market)

	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", loginHandler.ServeHTTP)
			r.Post("/register", registerHandler.ServeHTTP)

			r.Get("/orders", listOrdersHandler.ServeHTTP)
			r.Post("/orders", addOrdersHandler.ServeHTTP)
		})
	})
	return router
}
