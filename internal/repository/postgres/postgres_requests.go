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
	var (
		query = fmt.Sprintf("UPDATE %s SET accrual = $1 "+
			"where orderid = $2 and userid = $3",
			domain.TableOrders)
	)
	rows, err := pg.db.ExecContext(ctx, query, withdraw.Sum, withdraw.OrderID, user.ID)
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

func (pg *PostgresRepo) GetUserWithdrawals(ctx context.Context, user *users.User) (wd []withdraws.Withdraw, err error) {
	//TODO implement me
	panic("implement me")
}
