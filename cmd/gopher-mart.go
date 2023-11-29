package main

import (
	"gopher-mart/internal/config"
	"gopher-mart/internal/repository/postgres"
	httpserver "gopher-mart/internal/transport/http"
	gophermart "gopher-mart/internal/usecase/gopher-mart"
)

func main() {

	err := config.LoadEnvfileConfig()
	if err != nil {
		return
	}
	conf := config.LoadConfig()
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
