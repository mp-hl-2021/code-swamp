package codesnippet

import (
	"errors"
	"fmt"
	"github.com/mp-hl-2021/code-swamp/internal/domain/codesnippet"
	"github.com/mp-hl-2021/code-swamp/internal/usecases/account"
	"strings"
	"time"
)

var (
	ErrorUnsupportedLanguage = errors.New("unsupported language")
)

type CodeCheckResult struct {
	correct bool
	msg     string
}

type Interface interface {
	GetMySnippetIds(a account.Account) ([]uint, error)
	CreateSnippet(a *account.Account, code string, lang string, lifetime time.Duration) (uint, error)
	GetSnippetById(uint) (codesnippet.CodeSnippet, error)
	CheckCode(code string, lang string) (CodeCheckResult, error)
	SetIncorrectCode(r CodeCheckResult, err error, sid uint) error
}

type UseCases struct {
	CodeSnippetStorage codesnippet.Interface
	CodeCheckChannel   chan<- CodeCheckResult
}

func (u *UseCases) CheckCode(code string, lang string) (CodeCheckResult, error) {
	// что-то гениальное
}

func (u *UseCases) SetIncorrectCode(r CodeCheckResult, err error, sid uint) error {
	if err != nil {
		u.CodeSnippetStorage.SetCodeStatus(sid, false, err.Error())
	} else if !r.correct {
		u.CodeSnippetStorage.SetCodeStatus(sid, false, r.msg)
	} else {
		u.CodeSnippetStorage.SetCodeStatus(sid, true, "")
	}
	return nil
}

func (u *UseCases) GetMySnippetIds(a account.Account) ([]uint, error) {
	if err := u.CodeSnippetStorage.DeleteExpiredSnippets(); err != nil {
		return []uint{}, err
	}
	fmt.Printf("GetMySnippetIds: %d\n", a.Id)
	return u.CodeSnippetStorage.GetMyCodeSnippetIds(a.Id)
}

func (u *UseCases) CreateSnippet(a *account.Account, code string, lang string, lifetime time.Duration) (uint, error) {
	shortenedCode := code
	if len(code) > 10 {
		shortenedCode = code[:10] + "..."
	}
	fmt.Printf("CreateSnippet: %s\n", shortenedCode)
	if lang != "" {
		if err := validateLanguage(lang); err != nil {
			return 0, err
		}
	}
	s := codesnippet.CodeSnippet{
		Code:      code,
		Lang:      lang,
		IsChecked: false,
		IsCorrect: false,
		Lifetime:  lifetime,
	}
	var sid uint
	var err error
	if a == nil {
		sid, err = u.CodeSnippetStorage.CreateCodeSnippet(s)
		if err != nil {
			return 0, err
		}
	} else {
		sid, err = u.CodeSnippetStorage.CreateCodeSnippetWithUser(s, a.Id)
		if err != nil {
			return 0, err
		}
	}
	go func() {
		r, err := u.CheckCode(code, lang)
		if err != nil || r.correct {
			if err := u.SetIncorrectCode(r, err, sid); err != nil {
				fmt.Printf("Falied to set code status: %s\n", err)
			}
		}
	}()
	return sid, nil
}

func (u *UseCases) GetSnippetById(id uint) (codesnippet.CodeSnippet, error) {
	if err := u.CodeSnippetStorage.DeleteExpiredSnippets(); err != nil {
		return codesnippet.CodeSnippet{}, err
	}
	fmt.Printf("GetSnippetById: %d\n", id)
	return u.CodeSnippetStorage.GetCodeSnippetById(id)
}

var SupportedLanguages = []string{"Python", "JavaScript", "Java", "Kotlin", "C#", "C", "C++", "PHP", "Swift", "Go", "Rust", "PETOOH"}

func validateLanguage(lang string) error {
	for _, l := range SupportedLanguages {
		if strings.ToLower(lang) == strings.ToLower(l) {
			return nil
		}
	}
	return ErrorUnsupportedLanguage
}
