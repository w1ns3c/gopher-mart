package accruals

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/domain/accruals"
	locerrors "gopher-mart/internal/domain/errors"
	"gopher-mart/internal/repository"
	"net/http"
	"time"
)

type AccrualsInf interface {
	CheckAccruals(ctx context.Context)
	SaveAccruals(ctx context.Context, ch chan *accruals.Accrual)
}

type Usecase struct {
	// Accrual system addr
	Addr string

	// worker pool
	WorkersCount int
	Attempts     int

	// repository
	repo repository.AccrualsRepoInf

	// TickerTimer
	timer time.Duration
}

func (u *Usecase) CheckAccruals(ctx context.Context) {
	ticker := time.NewTicker(u.timer)

	for range ticker.C {
		log.Info().Msg("Requesting accruals")
		ordersCh := make(chan string)

		// get orders ID from DB
		go u.GetProcessingOrders(ctx, ordersCh)
		accrualsCh := u.GetAccrualsFromRemote(ctx, ordersCh)
		go u.SaveAccruals(ctx, accrualsCh)
	}

}

func (u *Usecase) SaveAccruals(ctx context.Context, ch chan *accruals.Accrual) {
	for accrual := range ch {
		err := u.repo.UpdateAccrual(ctx, accrual)
		if err != nil {
			log.Error().Err(err).Send()
		}
	}
	log.Info().Str("status", "closed").Msg("DB")
}

func (u *Usecase) GetProcessingOrders(ctx context.Context, ordersCh chan string) error {
	ordersSl, err := u.repo.GetProccessingOrders(ctx)
	if err != nil {
		return err
	}
	go func(ordersSl []string) {
		for _, order := range ordersSl {
			ordersCh <- order
		}
	}(ordersSl)
	return nil
}

func (u *Usecase) GetAccrualsFromRemote(ctx context.Context,
	ordersID chan string) chan *accruals.Accrual {

	accrualsCh := make(chan *accruals.Accrual)

	for id := 0; id < u.WorkersCount; id++ {
		go u.accrualWorker(ctx, id, ordersID, accrualsCh)
	}
	return accrualsCh
}

func (u *Usecase) accrualWorker(ctx context.Context, workerID int,
	ordersID chan string, accrualsCh chan *accruals.Accrual) {

	for orderID := range ordersID {
		retry(func() error {
			// TODO Add retryableErrors
			return u.createRequest(ctx, orderID, accrualsCh)
		},
			u.Attempts, locerrors.ErrTooManyRequests)
	}

}

func (u *Usecase) createRequest(ctx context.Context, orderID string,
	accrualsCh chan *accruals.Accrual) error {

	url := fmt.Sprintf("%s/%s", u.Addr, orderID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return locerrors.ErrTooManyRequests
	}
	newAccrual := new(accruals.Accrual)
	err = json.NewDecoder(resp.Body).Decode(&newAccrual)
	if err != nil {
		return err
	}

	// all successful
	accrualsCh <- newAccrual
	return nil
}

func retry(f func() error, attempts int, retryableErrors ...error) {
	for i := 0; i < attempts; i++ {
		err := f()
		if err != nil {
			for _, e := range retryableErrors {
				if errors.Is(err, e) {
					time.Sleep(time.Second * 90)
					continue
				}
			}
			break
		}
	}
}
