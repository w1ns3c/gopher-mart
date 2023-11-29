package main

import (
	"gopher-mart/internal/config"
	"gopher-mart/internal/repository/postgres"
	httpserver "gopher-mart/internal/transport/http"
	gophermart "gopher-mart/internal/usecase/gopher-mart"
)

func main() {

	conf, err := config.LoadConfig()
	if err != nil {
		// TODO log err
		return
	}

	repo := postgres.NewRepository(conf.DBurl)
	market := gophermart.NewGophermart(
		gophermart.WithRepo(repo),
	)
	router := httpserver.NewRouter(market)

	srv, err := httpserver.NewHTTPServer(conf.SrvAddr, router)
	if err != nil {
		return
	}
	srv.ListenAndServe()

}
