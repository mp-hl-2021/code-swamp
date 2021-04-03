package _interface

import (
	"github.com/mp-hl-2021/code-swamp/internal/domain/repository"
	"sync"
)

type Memory struct {
	accountsById    map[uint]repository.Account
	accountsByLogin map[string]repository.Account
	nextId          uint
	mu              *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		accountsById:    make(map[uint]repository.Account),
		accountsByLogin: make(map[string]repository.Account),
		mu:              &sync.Mutex{},
	}
}

func (m *Memory) CreateAccount(cred repository.Credentials) (repository.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.accountsByLogin[cred.Login]; ok {
		return repository.Account{}, repository.ErrAlreadyExist
	}
	a := repository.Account {
		Id: m.nextId,
		Credentials: cred,
	}
	m.accountsById[a.Id] = a
	m.accountsByLogin[a.Login] = a
	m.nextId++
	return a, nil
}

func (m *Memory) GetAccountById(id uint) (repository.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsById[id]
	if !ok {
		return a, repository.ErrNotFound
	}
	return a, nil
}

func (m *Memory) GetAccountByLogin(login string) (repository.Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.accountsByLogin[login]
	if !ok {
		return a, repository.ErrNotFound
	}
	return a, nil
}