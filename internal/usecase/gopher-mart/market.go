package gophermart

import (
	"context"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/config"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	"gopher-mart/internal/repository"
	"gopher-mart/internal/repository/postgres"
	"gopher-mart/internal/usecase/cookies"
	ordersUsecase "gopher-mart/internal/usecase/orders"
	usersUsecase "gopher-mart/internal/usecase/users"
	"net/http"
	"time"
)

type MarketUsecaseInf interface {
	usersUsecase.UserUsecaseInf
	//usersUsecase.UserBalanceUsecase
	//usersUsecase.UserContextUsecase

	ordersUsecase.OrdersUsecaseInf
	ordersUsecase.OrderValidator
	cookies.CookiesUsecae

	WithdrawUserBonuses(ctx context.Context, user *users.User, wd *withdraws.Withdraw) error
}

type GopherMart struct {
	// optional params
	Secret         string
	CookieName     string
	CookieLifetime time.Duration

	// important params
	AccrualSystemHost string
	dbURL             string

	cookies *cookies.Usecase
	users   *usersUsecase.Usecase
	orders  *ordersUsecase.Usecase
	repo    repository.Repository

	ctx context.Context
}

// TODO Choose constructor
func NewGmartWithConfig(config *config.Config) (mart *GopherMart, err error) {
	mart = &GopherMart{
		// optional params
		Secret:         config.Secret,
		CookieName:     config.CookieName,
		CookieLifetime: config.CookieHoursLifeTime,

		// important params
		AccrualSystemHost: config.RemoteServiceAddr,
		dbURL:             config.DBurl,
	}

	// initialize repo
	repo, err := postgres.NewRepository(mart.dbURL, context.TODO())
	if err != nil {
		return nil, err
	}
	log.Info().Str("dbUrl", mart.dbURL).Msg("DB connected")
	mart.repo = repo

	// initialize usecases
	mart.users = usersUsecase.NewUsecaseWith(
		usersUsecase.WithRepo(mart.repo),
		usersUsecase.WithSecret(mart.Secret),
		usersUsecase.WithCookieName(mart.CookieName),
		usersUsecase.WithCookieLifetime(mart.CookieLifetime),
	)

	mart.orders = ordersUsecase.NewUsecaseWith(
		ordersUsecase.WithRepo(mart.repo),
	)
	mart.cookies = cookies.NewUsecaseWith(
		cookies.WithRepo(mart.repo),
		cookies.WithSecret(mart.Secret),
	)

	return mart, nil
}

// TODO Delete ??
// ------------------------------------------
type MartOptions func(mart *GopherMart)

func NewGophermart(options ...MartOptions) *GopherMart {
	market := new(GopherMart)
	for _, option := range options {
		option(market)
	}
	return market
}

func WithConfig(config *config.Config) func(mart *GopherMart) {
	return func(mart *GopherMart) {
		// optional params
		mart.Secret = config.Secret
		mart.CookieName = config.CookieName
		mart.CookieLifetime = config.CookieHoursLifeTime

		// important params
		mart.AccrualSystemHost = config.RemoteServiceAddr
		mart.dbURL = config.DBurl

	}
}

func WithContext(ctx context.Context) func(mart *GopherMart) {
	return func(mart *GopherMart) {
		mart.ctx = ctx
	}
}

func WithSecret(secret string) func(mart *GopherMart) {
	return func(mart *GopherMart) {
		mart.Secret = secret
	}
}

func WithCookieName(cookieName string) func(mart *GopherMart) {
	return func(mart *GopherMart) {
		mart.CookieName = cookieName
	}
}

func WithCookieLifetime(lifetime time.Duration) func(mart *GopherMart) {
	return func(mart *GopherMart) {
		mart.CookieLifetime = lifetime
	}
}
func WithRepo(repo repository.Repository) func(mart *GopherMart) {
	return func(mart *GopherMart) {
		mart.repo = repo
	}
}

func InitUsecases() func(mart *GopherMart) {
	return func(mart *GopherMart) {
		if mart.repo == nil {
			log.Fatal().Err(errors.ErrRepoNotInit).Send()
		}
		if mart.Secret == "" || mart.CookieName == "" ||
			mart.AccrualSystemHost == "" {
			log.Fatal().Err(errors.ErrGophermart).Send()
		}

		// initialize usecases
		mart.users = usersUsecase.NewUsecaseWith(
			usersUsecase.WithRepo(mart.repo),
			usersUsecase.WithSecret(mart.Secret),
			usersUsecase.WithCookieName(mart.CookieName),
			usersUsecase.WithCookieLifetime(mart.CookieLifetime),
		)

		mart.orders = ordersUsecase.NewUsecaseWith(
			ordersUsecase.WithRepo(mart.repo),
		)

		mart.cookies = cookies.NewUsecaseWith(
			cookies.WithRepo(mart.repo),
			cookies.WithSecret(mart.Secret),
		)

	}
}

// ------------------------------------------

func (g *GopherMart) LoginUser(ctx context.Context, user *users.User) (cookie *http.Cookie, err error) {
	return g.users.LoginUser(ctx, user)
}

func (g *GopherMart) RegisterUser(ctx context.Context, user *users.User) error {
	return g.users.RegisterUser(ctx, user)
}

func (g *GopherMart) GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error) {
	return g.users.GetUserWithdrawals(ctx, user)
}

func (g *GopherMart) CheckBalance(ctx context.Context, user *users.User) (balance *users.Balance, err error) {
	return g.users.CheckBalance(ctx, user)
}

func (g *GopherMart) CheckUserInContext(ctx context.Context) (user *users.User, err error) {
	return g.users.CheckUserInContext(ctx)
}

func (g *GopherMart) ListOrders(ctx context.Context, user *users.User) (orders []orders.Order, err error) {
	return g.orders.ListOrders(ctx, user)
}

func (g *GopherMart) AddOrder(ctx context.Context, user *users.User, orderNumber string) error {
	return g.orders.AddOrder(ctx, user, orderNumber)
}

func (g *GopherMart) WithdrawUserBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error {
	/*
		1. Request user balance
		2. Compare balance with Sum (withdraw request)
		3. Save to DB
		3.1. Save: new fields (withdraws) for orders
		3.2. Save: new balance and withdraws for user (balance table)
	*/

	// request balance
	balance, err := g.users.CheckBalance(ctx, user)
	if err != nil {
		return err
	}

	// compare balances
	if balance.Current < withdraw.Sum {
		log.Warn().Str("user", user.Login).Send()
		return errors.ErrNotEnoughBonuses
	}

	// 3.1 save withdraw for order
	err = g.repo.WithdrawBonuses(ctx, user, withdraw)
	if err != nil {
		//log.Error().Err(err).Send()
		return errors.ErrOrderWrongFormat
	}

	// 3.2 change balance
	balance.Current -= withdraw.Sum
	balance.WithdrawsSum += withdraw.Sum

	return g.repo.UpdateBalance(ctx, user, balance)
}

func (g *GopherMart) ValidateOrderFormat(orderNumber string) bool {
	return g.orders.ValidateOrderFormat(orderNumber)
}

func (g *GopherMart) ValidateCookie(ctx context.Context, cookie *http.Cookie) (user *users.User, err error) {
	return g.cookies.ValidateCookie(ctx, cookie)
}

func (g *GopherMart) UpdateBalance(ctx context.Context, user *users.User, balance *users.Balance) error {
	return g.users.UpdateBalance(ctx, user, balance)
}
