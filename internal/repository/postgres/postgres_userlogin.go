package postgres

import (
	"context"
	"fmt"
	"gopher-mart/internal/domain"
	"gopher-mart/internal/domain/errors"
	"gopher-mart/internal/domain/users"
	"time"
)

func (pg *PostgresRepo) LoginUser(ctx context.Context, user *users.User) (userID string, userHash string, err error) {
	var (
		query = fmt.Sprintf("SELECT userid, login, hash"+
			" FROM %s where login=$1;", domain.TableUsers)
	)
	rows, err := pg.db.QueryContext(ctx, query, user.Login)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	result := make([]users.User, 0)

	for rows.Next() {
		var (
			id          string
			login, hash string
		)

		if err := rows.Scan(&id, &login, &hash); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			//log.Fatal(err)
			return "", "", err
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
		return "", "", rerr
	}
	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return "", "", rerr
	}

	if len(result) != 1 {
		return "", "", errors.ErrWrongResultValues
	}

	return result[0].ID, result[0].Hash, nil
}

func (pg *PostgresRepo) RegisterUser(ctx context.Context, user *users.User) error {
	var (
		query1 = fmt.Sprintf("insert into %s (userid, login, hash) "+
			"values ($1, $2, $3);", domain.TableUsers)
		query2 = fmt.Sprintf("insert into %s (userid, balance, withdraw) "+
			"values ($1,0,0);", domain.TableBalance)
	)
	ctx, cancel := context.WithTimeout(ctx, time.Second*8)
	defer cancel()
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query1, user.ID, user.Login, user.Hash)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, query2, user.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (pg *PostgresRepo) CheckUserExist(ctx context.Context, userID string) (user *users.User, err error) {
	var (
		query = fmt.Sprintf("SELECT userid, login "+
			"FROM %s where userid=$1;", domain.TableUsers)
	)
	rows, err := pg.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]users.User, 0)

	for rows.Next() {
		var (
			id    string
			login string
		)

		if err := rows.Scan(&id, &login); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			//log.Fatal(err)
			return nil, err
		}
		result = append(result, users.User{
			ID:    id,
			Login: login,
		})
	}
	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
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
