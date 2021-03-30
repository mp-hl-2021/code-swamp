package accountstorage

import (
	"sync"
)

type Memory struct {
	accountsById    map[uint]Account
	accountsByLogin map[string]Account
	nextId          uint
	mu              *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		accountsById:    make(map[uint]Account),
		accountsByLogin: make(map[string]Account),
		mu:              &sync.Mutex{},
	}
}

func (m *Memory) CreateAccount(cred Credentials) (Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.accountsByLogin[cred.Login]; ok {
		return Account{}, ErrAlreadyExist
	}
	a := Account{
		Id: m.nextId,
		Credentials: cred,
	}
	m.accountsById[a.Id] = a
	m.accountsByLogin[a.Login] = a
	m.nextId++
	return a, nil
}

func (m *Memory) GetAccountById(id uint) (Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsById[id]
	if !ok {
		return a, ErrNotFound
	}
	return a, nil
}

func (m *Memory) GetAccountByLogin(login string) (Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsByLogin[login]
	if !ok {
		return a, ErrNotFound
	}
	return a, nil
}