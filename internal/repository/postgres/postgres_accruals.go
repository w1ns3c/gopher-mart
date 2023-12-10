package postgres

import (
	"context"
	"fmt"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/accruals"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/orders"
)

func (pg *PostgresRepo) UpdateAccrual(ctx context.Context, accrual *accruals.Accrual) error {
	var (
		query = fmt.Sprintf("UPDATE %s SET status=$1, accrual=$2 "+
			"where orderid=$3 and status<>$1", domain.TableOrders)
	)

	rows, err := pg.db.ExecContext(ctx, query, accrual.Status, accrual.Accrual, accrual.Order)
	if err != nil {
		return err
	}
	count, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.ErrAlreadyUpdated
	}
	return nil
}

func (pg *PostgresRepo) GetProccessingOrders(ctx context.Context) (ordersID []string, err error) {
	var (
		query = fmt.Sprintf("SELECT orderid FROM %s "+
			"WHERE status<>$1 and "+
			"status <>$2 "+
			"ORDER BY upload_date "+
			"LIMIT 1000", domain.TableOrders)
	)

	rows, err := pg.db.QueryContext(ctx, query, orders.StatusInvalid, orders.StatusDone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]string, 0)
	for rows.Next() {
		var orderID string
		err = rows.Scan(&orderID)
		if err != nil {
			return nil, err
		}
		result = append(result, orderID)
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

func (pg *PostgresRepo) GetUserByOrderID(ctx context.Context, orderID string) (userID string, err error) {
	var (
		query = fmt.Sprintf("SELECT userid "+
			"FROM %s where orderID=$1;", domain.TableOrders)
	)
	rows, err := pg.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	result := make([]string, 0)

	for rows.Next() {
		var (
			id string
		)

		if err := rows.Scan(&id); err != nil {
			return "", err
		}
		result = append(result, id)
	}
	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	rerr := rows.Close()
	if rerr != nil {
		return "", rerr
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		return "", err
	}

	if len(result) != 1 {
		return "", errors.ErrWrongResultValues
	}

	return result[0], nil
}
