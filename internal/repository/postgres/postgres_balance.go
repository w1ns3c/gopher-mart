package postgres

import (
	"context"
	"fmt"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
)

func (pg *PostgresRepo) CheckBalance(ctx context.Context, user *users.User) (balance *users.Balance, err error) {
	var (
		query = fmt.Sprintf("SELECT balance, withdraw FROM %s where userid=$1;", domain.TableBalance)
	)

	rows, err := pg.db.QueryContext(ctx, query, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]users.Balance, 0)
	for rows.Next() {
		var (
			balance users.Balance
		)

		err := rows.Scan(&balance.Current, &balance.WithdrawsSum)
		if err != nil {
			return nil, err
		}
		result = append(result, balance)
	}

	rerr := rows.Close()
	if rerr != nil {
		return nil, rerr
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return nil, rerr
	}

	if len(result) != 1 {
		return nil, errors.ErrWrongResultValues
	}

	return &result[0], nil
}

func (pg *PostgresRepo) UpdateBalance(ctx context.Context, user *users.User, balance *users.Balance) error {
	var (
		query = fmt.Sprintf("UPDATE %s SET balance=$1, withdraw=$2 WHERE userid=$3", domain.TableBalance)
	)

	_, err := pg.db.ExecContext(ctx, query, balance.Current, balance.WithdrawsSum, user.ID)
	return err
}
