package domain

import (
	"database/sql"
	"fmt"

	"github.com/alflanagan/banking-lib/errs"
	"github.com/alflanagan/banking-lib/logger"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	FindBy(username string, password string) (*Login, *errs.AppError)
	GenerateRefreshTokenAndStore(authToken AuthToken) (string, *errs.AppError)
}

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (d AuthRepositoryDb) GenerateRefreshTokenAndStore(authToken AuthToken) (string, *errs.AppError) {
	var appErr *errs.AppError
	var refreshToken string

	// generate the refresh token
	if refreshToken, appErr = authToken.newRefreshToken(); appErr != nil {
		return "", appErr
	}
	// store it in the store
	sqlInsert := "INSERT INTO banking.refresh_token_store(refresh_token) VALUES (?)"
	if _, err := d.client.Exec(sqlInsert, refreshToken); err != nil {
		logger.Error("Unexpected db error: " + err.Error())
		return "", errs.NewUnexpectedError("unexpected database error")
	}
	return refreshToken, nil
}

func (d AuthRepositoryDb) FindBy(username, password string) (*Login, *errs.AppError) {
	var login Login
	sqlVerify := `SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
                  LEFT JOIN accounts a ON a.customer_id = u.customer_id
                WHERE username = ? and password = ?
                GROUP BY a.customer_id`
	err := d.client.Get(&login, sqlVerify, username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("invalid credentials")
		} else {
			logger.Info(fmt.Sprintf("Error while verifying login request from database: %s", err.Error()))
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &login, nil
}

func NewAuthRepository(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}
