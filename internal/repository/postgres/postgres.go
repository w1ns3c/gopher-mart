package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/domain"
	"strings"
)

type PostgresRepo struct {
	db  *sql.DB
	url string
}

func NewRepository(url string) *PostgresRepo {
	if !strings.Contains(url, "postgres://") {
		url = "postgres://" + url
	}
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Error().Err(err).Send()
		return nil
	}
	return &PostgresRepo{
		db:  db,
		url: url,
	}
}

func (pg *PostgresRepo) CheckConnection(ctx context.Context) error {
	var err error
	log.Info().Str("db_url", pg.url).Send()
	pg.db, err = sql.Open("pgx", pg.url)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}
	err = pg.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresRepo) Init(ctx context.Context) error {
	if pg.db == nil {
		return fmt.Errorf("db not connected")
	}
	err := pg.CheckConnection(ctx)
	if err != nil {
		return err
	}

	queryTb1 := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
		userid integer NOT NULL UNIQUE,
		login varchar NOT NULL UNIQUE,
		hash varchar NOT NULL,
		CONSTRAINT users_pk PRIMARY KEY (userid),
		CONSTRAINT users_fk FOREIGN KEY (userid) REFERENCES public.users(userid)
	);`, domain.TableUsers)

	queryTb2 := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		userid integer NOT NULL UNIQUE,
		balance integer NULL,
		withdraw integer NULL,
		CONSTRAINT balance_pk PRIMARY KEY (userid),
		CONSTRAINT balance_fk FOREIGN KEY (userid) REFERENCES public.balance(userid)
	);`, domain.TableOrders)

	queryTb3 := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		orderid integer NOT NULL,
		userid integer NOT NULL,
		status varchar NOT NULL,
		accrual integer NOT NULL,
		upload_date varchar NULL,
		CONSTRAINT orders_pk PRIMARY KEY (orderid),
		CONSTRAINT orders_fk FOREIGN KEY (orderid) REFERENCES public.orders(orderid)
	);`, domain.TableBalance)

	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, queryTb1)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, queryTb2)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, queryTb3)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
