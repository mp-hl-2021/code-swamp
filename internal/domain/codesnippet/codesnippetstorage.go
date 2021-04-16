package codesnippet

import (
	"time"
)

type CodeSnippet struct {
	Code     string
	Lang     string
	Lifetime time.Duration
}

type Interface interface {
	CreateCodeSnippet(s CodeSnippet) (uint, error)
	CreateCodeSnippetWithUser(s CodeSnippet, uid uint) (uint, error)
	GetCodeSnippetById(sid uint) (CodeSnippet, error)
	GetMyCodeSnippetIds(uid uint) ([]uint, error)
	DeleteExpiredSnippets() error
}
