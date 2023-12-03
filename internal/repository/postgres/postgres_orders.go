package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/orders"
	"gopher-mart/internal/domain/users"
	"time"
)

func (pg *PostgresRepo) AddOrder(ctx context.Context, user *users.User, orderNumber string) error {
	var (
		queryOrders = fmt.Sprintf("INSERT INTO %s (orderid, userid, status, upload_date) "+
			"values ($1, $2, $3, $4)", domain.TableOrders)
		queryWithdraws = fmt.Sprintf("INSERT INTO %s (orderid, userid) "+
			"values ($1, $2)", domain.TableWithdraws)
	)

	now := time.Now()
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, queryOrders, orderNumber, user.ID,
		orders.StatusNew, now)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.ExecContext(ctx, queryWithdraws, orderNumber, user.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pg *PostgresRepo) GetUserOrders(ctx context.Context, user *users.User) (result []orders.Order, err error) {
	var (
		query = fmt.Sprintf("SELECT orderid, status, accrual, upload_date "+
			"FROM %s WHERE userid=$1 ORDER BY upload_date DESC;", domain.TableOrders)
	)

	rows, err := pg.db.QueryContext(ctx, query, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result = make([]orders.Order, 0)
	for rows.Next() {
		var (
			order    orders.Order
			date     sql.NullString
			status   sql.NullString
			cashback sql.NullInt32
		)

		err = rows.Scan(&order.ID, &status, &cashback, &date)
		if err != nil {
			return nil, err
		}

		if date.Valid {
			t, err := time.Parse(time.RFC3339, date.String)
			if err != nil {
				return nil, err
			}
			order.Date = t
		}
		order.Status = orders.OrderStatus(status.String)
		// TODO maybe, should use another type (conversion int32 in uint64)
		order.Cashback = uint64(cashback.Int32)

		result = append(result, order)
	}

	rerr := rows.Close()
	if rerr != nil {
		return nil, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (pg *PostgresRepo) CheckOrder(ctx context.Context, orderNumber string) (orderid, userid string, err error) {
	var (
		query = fmt.Sprintf("SELECT orderid, userid FROM %s "+
			"where orderid=$1;", domain.TableOrders)
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
