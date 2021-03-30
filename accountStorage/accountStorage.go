package accountstorage

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
)

type Account struct {
	Id uint
	Credentials
}

type Credentials struct {
	Login    string
	Password string
}

type Interface interface {
	CreateAccount(cred Credentials) (Account, error)
	GetAccountById(id uint) (Account, error)
	GetAccountByLogin(login string) (Account, error)
}