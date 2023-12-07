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
	"os"
	"os/signal"
	"syscall"
)

func main() {

	conf := config.LoadConfig()
	fmt.Println(conf)

	err := logger.Initialize(conf.LogLevel)
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	// initialise all context, service, repo and transport server
	// init context
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer cancel()

	// init repository
	repo, err := postgres.NewRepository(conf.DBurl, ctx)
	if err != nil {
		log.Error().Err(err).Msg("Repo init: ")
		return
	}

	log.Info().Msg(fmt.Sprint(conf))
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
	if err != nil {
		log.Error().Err(err).Send()
	}
	err = market.Close()
	if err != nil {
		log.Error().Err(err).Send()
	}

	<-ctx.Done()

	log.Info().Msg("Gophermart stopped")

}
