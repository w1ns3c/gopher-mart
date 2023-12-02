package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"gopher-mart/internal/domain"
	"strings"
)

type PostgresRepo struct {
	db  *sql.DB
	url string
}

func NewRepository(dbURL string, ctx context.Context) (repo *PostgresRepo, err error) {
	if !strings.Contains(dbURL, "postgres://") {
		dbURL = "postgres://" + dbURL
	}
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}
	repo = &PostgresRepo{
		db:  db,
		url: dbURL,
	}
	return repo, repo.Init(ctx)
}

func (pg *PostgresRepo) CheckConnection() error {
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
	err := pg.CheckConnection()
	if err != nil {
		return err
	}

	queryTb1 := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
		userid varchar primary KEY,
		login varchar NOT NULL UNIQUE,
		hash varchar NOT NULL,
		CONSTRAINT users_fk FOREIGN KEY (userid) REFERENCES public.users(userid)
	);`, domain.TableUsers)

	queryTb2 := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		userid varchar primary KEY,
		balance integer NULL,
		withdraw integer NULL,
		CONSTRAINT balance_fk FOREIGN KEY (userid) REFERENCES public.balance(userid)
	);`, domain.TableBalance)

	queryTb3 := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		orderid varchar primary KEY,
		userid varchar NOT NULL,
		status varchar ,
		accrual integer,
		upload_date varchar,
		CONSTRAINT orders_fk FOREIGN KEY (orderid) REFERENCES public.orders(orderid)
	);`, domain.TableOrders)

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
	return tx.Commit()
}
