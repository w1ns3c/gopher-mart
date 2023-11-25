package http

import (
	"github.com/go-chi/chi/v5"
	"gopher-mart/internal/transport/http/handlers"
	"gopher-mart/internal/transport/http/middlewares"
	market "gopher-mart/internal/usecase/gopher-mart"
	"net/http"
)

func NewRouter(market market.MarketUsecase) http.Handler {
	// init middlewares
	authMid := middlewares.NewAuthMidleware(market)

	// init handlers
	// login handlers
	loginHandler := handlers.NewLoginHandler(market)
	registerHandler := handlers.NewRegisterHandler(market)

	// orders handlers
	listOrdersHandler := handlers.NewOrdersListHandler(market)
	addOrdersHandler := handlers.NewOrdersAddHandler(market)

	// balance handlers
	getBalanceHandler := handlers.NewBalanceHandler(market)
	balanceWithdraw := handlers.NewBalanceWithdrawHandler(market)

	// withdraws handler
	allUserWithdraws := handlers.NewWithdrawalsHandler(market)

	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/login", loginHandler.ServeHTTP)
			r.Post("/register", registerHandler.ServeHTTP)

			// authed api
			r.Use(authMid.AuthMiddleware)

			// orders api
			r.Get("/orders", listOrdersHandler.ServeHTTP)
			r.Post("/orders", addOrdersHandler.ServeHTTP)

			// balance api
			r.Get("/balance", getBalanceHandler.ServeHTTP)
			r.Post("/balance/withdraw", balanceWithdraw.ServeHTTP)

			// withdraws api
			r.Get("/withdrawals", allUserWithdraws.ServeHTTP)

		})
	})

	return router
}
