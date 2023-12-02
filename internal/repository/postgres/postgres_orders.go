package postgres

import (
	"context"
	"fmt"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
)

func (pg *PostgresRepo) AddOrder(ctx context.Context, user *users.User, orderNumber string) error {
	var (
		query = fmt.Sprintf("INSERT INTO %s (orderid, userid) values ($1, $2)", domain.TableOrders)
	)

	_, err := pg.db.ExecContext(ctx, query, orderNumber, user.ID)
	return err
}

func (pg *PostgresRepo) CheckOrder(ctx context.Context, orderNumber string) (orderid, userid string, err error) {
	var (
		query = fmt.Sprintf("SELECT orderid, userid FROM %s where orderid=$1;", domain.TableOrders)
	)

	rows, err := pg.db.QueryContext(ctx, query, orderNumber)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	type Result struct {
		Orderid string
		Userid  string
	}

	result := make([]Result, 0)
	for rows.Next() {
		var (
			orderInfo Result
		)

		err = rows.Scan(&orderInfo.Orderid, &orderInfo.Userid)
		if err != nil {
			return "", "", err
		}
		result = append(result, orderInfo)
	}

	rerr := rows.Close()
	if rerr != nil {
		return "", "", err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		return "", "", err
	}

	if len(result) != 1 {
		return "", "", errors.ErrWrongResultValues
	}

	return result[0].Orderid, result[0].Userid, nil
}
