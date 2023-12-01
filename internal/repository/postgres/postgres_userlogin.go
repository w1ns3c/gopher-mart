package postgres

import (
	"context"
	"fmt"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
)

func (pg *PostgresRepo) LoginUser(ctx context.Context, user *users.User) (userID uint64, userHash string, err error) {
	var (
		query = fmt.Sprintf("SELECT userid, login, hash FROM %s where login=$1;", domain.TableUsers)
	)
	rows, err := pg.db.QueryContext(ctx, query, user.Login)
	if err != nil {
		return 0, "", err
	}
	defer rows.Close()
	result := make([]users.User, 0)

	for rows.Next() {
		var (
			id          uint64
			login, hash string
		)

		if err := rows.Scan(&id, &login, &hash); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			//log.Fatal(err)
			return 0, "", err
		}
		result = append(result, users.User{
			ID:    id,
			Login: login,
			Hash:  hash,
		})
	}
	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	rerr := rows.Close()
	if rerr != nil {
		return 0, "", rerr
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return 0, "", rerr
	}

	if len(result) != 1 {
		return 0, "", errors.ErrWrongResultValues
	}

	return result[0].ID, result[0].Hash, nil
}

func (pg *PostgresRepo) RegisterUser(ctx context.Context, user *users.User) error {
	var (
		query = fmt.Sprintf("insert into %s (login, hash) values ($1, $2);", domain.TableUsers)
	)
	result, err := pg.db.ExecContext(ctx, query, user.Login, user.Hash)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected to affect 1 row, affected %d", rows)
	}
	return err
}
