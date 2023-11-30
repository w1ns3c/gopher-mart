package httpserver

import (
	"github.com/go-chi/chi/v5"
	"gopher-mart/internal/transport/http/handlers"
	"gopher-mart/internal/transport/http/middlewares"
	market "gopher-mart/internal/usecase/gopher-mart"
	"net/http"
)

func NewRouter(market market.MarketUsecaseInf) http.Handler {
	// init middlewares
	authMid := middlewares.NewAuthMidleware(market)
	//ddosMid := middlewares.NewDDOSMiddleware(market.GetMaxRequestsPerMinute())

	// init handlers
	// login handlers
	loginHandler := handlers.NewLoginHandler(market)
	registerHandler := handlers.NewRegisterHandler(market)

	// orders handlers
	listOrdersHandler := handlers.NewOrdersListHandler(market)
	addOrdersHandler := handlers.NewOrdersAddHandler(market)
	//orderStatusHandler := handlers.NewOrderStatusHandler(market)

	// balance handlers
	getBalanceHandler := handlers.NewBalanceHandler(market)
	balanceWithdraw := handlers.NewBalanceWithdrawHandler(market)

	// withdraws handler
	allUserWithdraws := handlers.NewWithdrawalsHandler(market)

	router := chi.NewRouter()

	// use logging middleware
	router.Use(middlewares.LoggingMiddleware)

	// use gzip compression
	router.Use(middlewares.GzipMiddleware)

	// routing
	router.Route("/api/user", func(r chi.Router) {
		r.Post("/login", loginHandler.ServeHTTP)
		r.Post("/register", registerHandler.ServeHTTP)

		// authed api
		r.Route("/", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.NotFound(http.NotFound)
			})
			r.Group(func(r chi.Router) {
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
	})

	//r.Route("/orders", func(r chi.Router) {
	//	r.Use(ddosMid.DDOSMiddleware)
	//	r.Get("/{number}", orderStatusHandler.ServeHTTP)
	//
	//})

	//})

	return http.Handler(router)
}
