package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/config"
	"gopher-mart/internal/logger"
	"gopher-mart/internal/repository/postgres"
	httpserver "gopher-mart/internal/transport/http"
	gophermart "gopher-mart/internal/usecase/gopher-mart"
)

func main() {

	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		log.Fatal().Err(err).Send()
		return
	}

	err = logger.Initialize(conf.LogLevel)
	if err != nil {
		fmt.Println(err)
		log.Fatal().Err(err).Send()
		return
	}

	// initialise all context, service, repo and transport server
	// init context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// init repository
	repo, err := postgres.NewRepository(conf.DBurl, ctx)
	if err != nil {
		fmt.Println(err)
		log.Error().Err(err).Msg("Repo init: ")
		return
	}

	fmt.Println(conf)
	// init usecases
	market := gophermart.NewGophermart(
		gophermart.WithRepo(repo),
		gophermart.WithConfig(conf),
		gophermart.InitUsecases(),
	)
	go market.CheckAccruals(ctx)

	router := httpserver.NewRouter(market)
	srv, err := httpserver.NewHTTPServer(conf.SrvAddr, router)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	// starting HTTP server
	err = srv.Run(ctx)
	log.Error().Err(err).Send()

}
