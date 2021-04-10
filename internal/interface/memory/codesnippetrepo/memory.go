package codesnippetrepo

import (
	"errors"
	"github.com/mp-hl-2021/code-swamp/internal/domain/codesnippet"
	"sync"
	"time"
)

var (
	ErrInvalidSnippedId = errors.New("no such snippet")
)

type SnippetInfo struct {
	cs         codesnippet.CodeSnippet
	uid        uint
	userExists bool
	exptime    time.Time
}

type Memory struct {
	snippetById map[uint]SnippetInfo
	nextId      uint
	mu          *sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		snippetById: make(map[uint]SnippetInfo),
		nextId:      0,
		mu:          &sync.Mutex{},
	}
}

func (m *Memory) CreateCodeSnippet(s codesnippet.CodeSnippet) (uint, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sid := m.nextId
	m.nextId += 1
	m.snippetById[sid] = SnippetInfo{
		cs:         s,
		uid:        0,
		userExists: false,
		exptime:    time.Now().Add(s.Lifetime),
	}
	return sid, nil
}

func (m *Memory) CreateCodeSnippetWithUser(s codesnippet.CodeSnippet, uid uint) (uint, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sid := m.nextId
	m.nextId += 1
	m.snippetById[sid] = SnippetInfo{
		cs:         s,
		uid:        uid,
		userExists: true,
		exptime:    time.Now().Add(s.Lifetime),
	}
	return sid, nil
}

func (m *Memory) GetCodeSnippetById(sid uint) (codesnippet.CodeSnippet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.snippetById[sid]
	if !ok {
		return codesnippet.CodeSnippet{}, ErrInvalidSnippedId
	}
	return s.cs, nil
}

func (m *Memory) GetMyCodeSnippetIds(uid uint) ([]uint, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var ids []uint
	for sid, i := range m.snippetById {
		if i.userExists && i.uid == uid {
			ids = append(ids, sid)
		}
	}
	return ids, nil
}

func (m *Memory) DeleteExpiredSnippets() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	validSnippets := make(map[uint]SnippetInfo)
	for sid, i := range m.snippetById {
		if i.exptime.After(time.Now()) {
			validSnippets[sid] = i
		}
	}
	m.snippetById = validSnippets
	return nil
}
