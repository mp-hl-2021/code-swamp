package codesnippetrepo

import (
	"errors"
	"github.com/mp-hl-2021/code-swamp/internal/domain/codesnippet"
	"sync"
)

var (
	ErrInvalidSnippedId = errors.New("no such snippet")
)

type Memory struct {
	snippetById       map[uint]codesnippet.CodeSnippet
	snippetIdsForUser map[uint][]uint
	nextId            uint
	mu                *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		snippetById:       make(map[uint]codesnippet.CodeSnippet),
		snippetIdsForUser: make(map[uint][]uint),
		nextId:            0,
		mu:                &sync.Mutex{},
	}
}

func (m *Memory) CreateCodeSnippet(s codesnippet.CodeSnippet) (uint, error) {
	sid := m.nextId
	m.nextId += 1
	m.snippetById[sid] = s
	return sid, nil
}

func (m *Memory) CreateCodeSnippetWithUser(s codesnippet.CodeSnippet, uid uint) (uint, error) {
	sid, err := m.CreateCodeSnippet(s)
	if err != nil {
		return 0, err
	}
	m.snippetIdsForUser[uid] = append(m.snippetIdsForUser[uid], sid)
	return sid, nil
}

func (m *Memory) GetCodeSnippetById(sid uint) (codesnippet.CodeSnippet, error) {
	s, ok := m.snippetById[sid]
	if !ok {
		return codesnippet.CodeSnippet{}, ErrInvalidSnippedId
	}
	return s, nil
}

func (m *Memory) GetMyCodeSnippetIds(uid uint) ([]uint, error) {
	ids, ok := m.snippetIdsForUser[uid]
	if !ok {
		return []uint{}, nil
	}
	return ids, nil
}