package postgres

import (
	"context"
	"fmt"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"gopher-mart/internal/domain/withdraws"
)

func (pg *PostgresRepo) WithdrawBonuses(ctx context.Context, user *users.User, withdraw *withdraws.Withdraw) error {
	//TODO implement me
	panic("implement me")
}

func (pg *PostgresRepo) GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error) {
	//TODO implement me
	panic("implement me")
}

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
