package mysql

import (
	"database/sql"
	"errors"
	"time"
	"unap-auth/domain/model"

	"golang.org/x/crypto/bcrypt"
)

type MySQLUserAccountRepository struct {
	DB *sql.DB
}

var (
	ErrUserNotFound    = errors.New("user account not found")
	ErrInvalidPassword = errors.New("username or password is incorrect")
	ErrInactiveAccount = errors.New("user account is inactive")
)

func (repo *MySQLUserAccountRepository) FindAndValidateUserAccount(userName string, password string) (*model.UserAccount, error) {
	query := `
			SELECT Id, UserName, Password, StartDate, CaducityDate, Status 
			FROM Core_UserAccounts 
			WHERE UserName = ?
			`
	row := repo.DB.QueryRow(query, userName)

	user := &model.UserAccount{}

	var startDateStr, caducityDateStr string

	err := row.Scan(&user.Id, &user.Username, &user.Password, &startDateStr, &caducityDateStr, &user.Status)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidPassword
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return nil, err
	}

	caducityDate, err := time.Parse("2006-01-02", caducityDateStr)
	if err != nil {
		return nil, err
	}

	if user.Status != 7 || startDate.After(time.Now()) || caducityDate.Before(time.Now()) {
		return nil, ErrInactiveAccount
	}

	return user, nil
}
