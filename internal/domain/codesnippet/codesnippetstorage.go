package codesnippet

import (
	"time"
)

type CodeSnippet struct {
	Code      string
	Lang      string
	IsChecked bool
	IsCorrect bool
	Message   string
	Lifetime  time.Duration
}

type Interface interface {
	CreateCodeSnippet(s CodeSnippet) (uint, error)
	CreateCodeSnippetWithUser(s CodeSnippet, uid uint) (uint, error)
	GetCodeSnippetById(sid uint) (CodeSnippet, error)
	GetMyCodeSnippetIds(uid uint) ([]uint, error)
	DeleteExpiredSnippets() error
	SetCodeStatus(sid uint, status bool, msg string) error
}
