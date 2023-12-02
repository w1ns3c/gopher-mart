package postgres

import (
	"context"
	"fmt"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
)

func (pg *PostgresRepo) CheckBalance(ctx context.Context, user *users.User) (curBalance, withDraws uint64, err error) {
	var (
		query = fmt.Sprintf("SELECT userid, balance, withdraw FROM %s where userid=$1;", domain.TableBalance)
	)

	rows, err := pg.db.QueryContext(ctx, query, user.ID)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()
	type Balance struct {
		ID        string
		Current   uint64
		Withdraws uint64
	}
	result := make([]Balance, 0)
	for rows.Next() {
		var (
			balance Balance
		)

		err := rows.Scan(&balance.ID, &balance.Current, &balance.Withdraws)
		if err != nil {
			return 0, 0, err
		}
		result = append(result, balance)
	}

	rerr := rows.Close()
	if rerr != nil {
		return 0, 0, rerr
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return 0, 0, rerr
	}

	if len(result) != 1 {
		return 0, 0, errors.ErrWrongResultValues
	}

	return result[0].Current, result[0].Withdraws, nil
}

func (pg *PostgresRepo) UpdateBalance(ctx context.Context, user *users.User, balance *users.Balance) error {
	var (
		query = fmt.Sprintf("UPDATE %s SET balance=$1, withdraw=$2 WHERE userid=$3", domain.TableBalance)
	)

	_, err := pg.db.ExecContext(ctx, query, balance.Current, balance.WithdrawsSum, user.ID)
	return err
}
