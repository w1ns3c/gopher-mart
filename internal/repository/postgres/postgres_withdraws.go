package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
	"time"
)

func (pg *PostgresRepo) WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error {
	var (
		query = fmt.Sprintf("UPDATE %s SET withdraw = $1, processed_at = $2 "+
			"where orderid = $3 and userid = $4 and withdraw is NULL",
			domain.TableWithdraws)
	)
	now := time.Now()
	rows, err := pg.db.ExecContext(ctx, query, withdraw.Sum, now, withdraw.OrderID, user.ID)
	if err != nil {
		return err
	}
	count, err := rows.RowsAffected()
	//log.Info().Int64("rows", count).Send()
	// order not updated
	if count != 1 {
		return errors.ErrOrderWrongFormat
	}
	return err
}

func (pg *PostgresRepo) GetUserWithdrawals(ctx context.Context,
	user *users.User) (result []withdraws.Withdraw, err error) {
	var (
		query = fmt.Sprintf("SELECT orderid, withdraw, processed_at "+
			"FROM %s WHERE userid=$1 ORDER BY processed_at DESC;", domain.TableWithdraws)
	)

	rows, err := pg.db.QueryContext(ctx, query, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result = make([]withdraws.Withdraw, 0)
	for rows.Next() {
		var (
			withdraw withdraws.Withdraw
			cashback sql.NullInt32
			date     sql.NullString
		)

		err = rows.Scan(&withdraw.OrderID, &cashback, &date)
		if err != nil {
			return nil, err
		}

		if date.Valid {
			t, err := time.Parse(time.RFC3339, date.String)
			if err != nil {
				return nil, err
			}
			withdraw.Date = t
		}

		// TODO maybe, should use another type (conversion int32 in uint64)
		withdraw.Sum = uint64(cashback.Int32)

		result = append(result, withdraw)
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
