package main

import (
	httpserver "gopher-mart/internal/transport/http"
	gophermart "gopher-mart/internal/usecase/gopher-mart"
)

func main() {

	address := "localhost:8000"
	market := gophermart.NewGophermart()
	router := httpserver.NewRouter(market)

	srv, err := httpserver.NewHTTPServer(address, router)
	if err != nil {
		return
	}
	srv.ListenAndServe()

}
