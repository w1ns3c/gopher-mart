package accruals

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/domain/accruals"
	locerrors "gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/repository"
	"net/http"
	"time"
)

type AccrualsInf interface {
	CheckAccruals(ctx context.Context)
	SaveAccruals(ctx context.Context, ch chan *accruals.Accrual)
	GetProcessingOrders(ctx context.Context, ordersCh chan string) error
}

type Usecase struct {
	// Accrual system addr
	Addr string

	// worker pool
	WorkersCount uint
	Attempts     uint

	// repository
	repo repository.AccrualsRepoInf

	// TickerTimer
	timer time.Duration
}

type AccrualsOptions func(usecase *Usecase)

func NewAccrualsWith(options ...AccrualsOptions) *Usecase {
	usecase := new(Usecase)
	for _, option := range options {
		option(usecase)
	}
	return usecase
}

func WithRepo(repo repository.AccrualsRepoInf) func(usecase *Usecase) {
	return func(usecase *Usecase) {
		usecase.repo = repo
	}
}

func WithAddr(addr string) func(u *Usecase) {
	return func(u *Usecase) {
		u.Addr = addr
	}
}

func WithWorkersCount(count uint) func(u *Usecase) {
	return func(u *Usecase) {
		u.WorkersCount = count
	}
}

func WithAttempts(attempts uint) func(u *Usecase) {
	return func(u *Usecase) {
		u.Attempts = attempts
	}
}

func WithTimer(timer time.Duration) func(u *Usecase) {
	return func(u *Usecase) {
		u.timer = timer
	}
}

func (u *Usecase) CheckAccruals(ctx context.Context) {
	ticker := time.NewTicker(u.timer)
	ordersCh := make(chan string)
	defer close(ordersCh)

	// start workers pool
	accrualsCh := u.GetAccrualsFromRemote(ctx, ordersCh)
	defer close(accrualsCh)

	for range ticker.C {
		log.Info().Msg("Requesting accruals")

		// get orders ID from DB
		go func() {
			err := u.GetProcessingOrders(ctx, ordersCh)
			if err != nil {
				log.Error().Err(err).Send()
			}

		}()
		go u.SaveAccruals(ctx, accrualsCh)
	}

}

func (u *Usecase) SaveAccruals(ctx context.Context, ch chan *accruals.Accrual) {
	for accrual := range ch {
		log.Info().Msg("Updating accruals")
		err := u.repo.UpdateAccrual(ctx, accrual)
		if err != nil {
			err = fmt.Errorf("can't update accrual: %v", err)
			log.Error().Err(err).Send()
			continue

		}
		userID, err := u.repo.GetUserByOrderID(ctx, accrual.Order)
		if err != nil {
			err = fmt.Errorf("can't get userID by orderID: %v", err)
			log.Error().Err(err).Send()
			continue
		}
		user := &users.User{ID: userID}
		balance, err := u.repo.CheckBalance(ctx, user)
		if err != nil {
			err = fmt.Errorf("can't get user balance: %v", err)
			log.Error().Err(err).Send()
			continue
		}

		balance.Current += accrual.Accrual
		err = u.repo.UpdateBalance(ctx, user, balance)
		if err != nil {
			err = fmt.Errorf("can't update balance with accrual: %v", err)
			log.Error().Err(err).Send()
			continue
		}

	}
	log.Info().Str("status", "closed").Msg("DB")
}

func (u *Usecase) GetProcessingOrders(ctx context.Context, ordersCh chan string) error {
	ordersSl, err := u.repo.GetProccessingOrders(ctx)
	if err != nil {
		return err
	}
	if len(ordersSl) == 0 {
		return nil
	}

	log.Info().Int("orders", len(ordersSl)).Send()
	go func(ordersSl []string) {
		for _, order := range ordersSl {
			ordersCh <- order
		}
	}(ordersSl)
	return nil
}

// GetAccrualsFromRemote start workers
func (u *Usecase) GetAccrualsFromRemote(ctx context.Context,
	ordersID chan string) chan *accruals.Accrual {

	accrualsCh := make(chan *accruals.Accrual)

	for id := uint(0); id < u.WorkersCount; id++ {
		go u.accrualWorker(ctx, id, ordersID, accrualsCh)
	}
	return accrualsCh
}

func (u *Usecase) accrualWorker(ctx context.Context, workerID uint,
	ordersID chan string, accrualsCh chan *accruals.Accrual) {

	for orderID := range ordersID {
		log.Info().Uint("id", workerID).Msg("worker get job")
		select {
		case <-ctx.Done():
			log.Warn().Uint("id", workerID).Msg("worker stopped")
			return
		default:
			retry(func() error {
				// TODO Add retryableErrors
				return u.createRequest(orderID, accrualsCh)
			},
				u.Attempts, u.timer, locerrors.ErrTooManyRequests, locerrors.ErrOrderNotRegisteredInRemote)
		}

	}

}

func (u *Usecase) createRequest(orderID string,
	accrualsCh chan *accruals.Accrual) error {

	url := fmt.Sprintf("%s/api/orders/%s", u.Addr, orderID)
	log.Info().Str("url", url).Msg("Accrual request")
	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(fmt.Errorf("http connection")).Send()
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("orderid=%s, status=%d, err=%v",
			orderID, resp.StatusCode, locerrors.ErrTooManyRequests)
	}
	if resp.StatusCode == http.StatusNoContent {
		return fmt.Errorf("orderid=%s, status=%d, err=%v",
			orderID, resp.StatusCode, locerrors.ErrTooManyRequests)
	}

	newAccrual := new(accruals.Accrual)
	err = json.NewDecoder(resp.Body).Decode(&newAccrual)
	if err != nil {
		log.Error().Err(fmt.Errorf("json error"))
		return err
	}
	data, err := json.Marshal(newAccrual)
	if err != nil {
		return err
	}

	log.Info().Bytes("accrual", data).Send()
	// all successful
	accrualsCh <- newAccrual
	return nil
}

func retry(f func() error, attempts uint, retryTime time.Duration, retryableErrors ...error) {
	for i := uint(0); i < attempts; i++ {
		err := f()
		if err != nil {
			log.Error().Str("key", "retry").Err(err).Send()
			for _, e := range retryableErrors {
				if errors.Is(err, e) {
					time.Sleep(time.Second * 60)
					continue
				}
			}
			break
		}
	}
}
