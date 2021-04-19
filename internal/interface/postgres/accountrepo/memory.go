package accountrepo

import (
	"database/sql"
	account "github.com/mp-hl-2021/code-swamp/internal/domain/account"
)

type Postgres struct {
	conn *sql.DB
}

func New(conn *sql.DB) *Postgres {
	return &Postgres{conn: conn}
}

const queryCreateAccount = `
	INSERT INTO accounts(
		login,
		password
	) VALUES ($1, $2)
	RETURNING id
`

func (p Postgres) CreateAccount(cred account.Credentials) (account.Account, error) {
	a := account.Account{Credentials : cred}
	row := p.conn.QueryRow(queryCreateAccount, cred.Login, cred.Password)
	err := row.Scan(&a.Id)
	if err != nil && err == sql.ErrNoRows { //todo
		return account.Account{}, err
	}
	return a, nil
}

const queryGetAccountById = `
	SELECT
		id,
		login,
		password
	FROM accounts 
	WHERE id = $1
`

func (p Postgres) GetAccountById(id uint) (account.Account, error) {
	a := account.Account{}
	row := p.conn.QueryRow(queryGetAccountById, id)
	err := row.Scan(&a.Id, &a.Login, &a.Password)
	if err != nil && err == sql.ErrNoRows { //todo
		return account.Account{}, err
	}
	return a, nil
}

const queryGetAccountByLogin = `
	SELECT
		id,
		login,
		password
	FROM accounts 
	WHERE login = $1
`

func (p Postgres) GetAccountByLogin(login string) (account.Account, error) {
	a := account.Account{}
	row := p.conn.QueryRow(queryGetAccountById, login)
	err := row.Scan(&a.Id, &a.Login, &a.Password)
	if err != nil && err == sql.ErrNoRows { //todo
		return account.Account{}, err
	}
	if err != nil { //todo
		return account.Account{}, err
	}
	return a, nil
}
