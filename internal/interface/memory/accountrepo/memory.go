package accountrepo

import (
	account "github.com/mp-hl-2021/code-swamp/internal/domain/account"
	"sync"
)

type Memory struct {
	accountsById    map[uint]account.Account
	accountsByLogin map[string]account.Account
	nextId          uint
	mu              *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		accountsById:    make(map[uint]account.Account),
		accountsByLogin: make(map[string]account.Account),
		mu:              &sync.Mutex{},
	}
}

func (m *Memory) CreateAccount(cred account.Credentials) (account.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.accountsByLogin[cred.Login]; ok {
		return account.Account{}, account.ErrAlreadyExist
	}
	a := account.Account {
		Id: m.nextId,
		Credentials: cred,
	}
	m.accountsById[a.Id] = a
	m.accountsByLogin[a.Login] = a
	m.nextId++
	return a, nil
}

func (m *Memory) GetAccountById(id uint) (account.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsById[id]
	if !ok {
		return a, account.ErrNotFound
	}
	return a, nil
}

func (m *Memory) GetAccountByLogin(login string) (account.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsByLogin[login]
	if !ok {
		return a, account.ErrNotFound
	}
	return a, nil
}