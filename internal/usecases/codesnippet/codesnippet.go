package codesnippet

import (
	"fmt"
	"github.com/mp-hl-2021/code-swamp/internal/domain/codesnippet"
	"github.com/mp-hl-2021/code-swamp/internal/usecases/account"
	"time"
)

type Interface interface {
	GetMySnippetIds(a account.Account) ([]uint, error)
	CreateSnippet(a *account.Account, code string, lang string, lifetime time.Duration) (uint, error)
	GetSnippetById(uint) (codesnippet.CodeSnippet, error)
}

type UseCases struct {
	CodeSnippetStorage codesnippet.Interface
}

func (u *UseCases) GetMySnippetIds(a account.Account) ([]uint, error) {
	if err := u.CodeSnippetStorage.DeleteExpiredSnippets(); err != nil {
		return []uint{}, err
	}
	fmt.Printf("GetMySnippetIds: %d\n", a.Id)
	return u.CodeSnippetStorage.GetMyCodeSnippetIds(a.Id)
}

func (u *UseCases) CreateSnippet(a *account.Account, code string, lang string, lifetime time.Duration) (uint, error) {
	if err := u.CodeSnippetStorage.DeleteExpiredSnippets(); err != nil {
		return 0, err
	}
	fmt.Printf("CreateSnippet: %d %s\n", a.Id, code)
	if lang != "" {
		if err := validateLanguage(lang); err != nil {
			return 0, err
		}
	}
	s := codesnippet.CodeSnippet{
		Code:     code,
		Lang:     lang,
		Lifetime: lifetime,
	}
	if a == nil {
		sid, err := u.CodeSnippetStorage.CreateCodeSnippet(s)
		if err != nil {
			return 0, err
		}
		return sid, nil
	} else {
		sid, err := u.CodeSnippetStorage.CreateCodeSnippetWithUser(s, a.Id)
		if err != nil {
			return 0, err
		}
		return sid, nil
	}
}

func (u *UseCases) GetSnippetById(id uint) (codesnippet.CodeSnippet, error) {
	if err := u.CodeSnippetStorage.DeleteExpiredSnippets(); err != nil {
		return codesnippet.CodeSnippet{}, err
	}
	fmt.Printf("GetSnippetById: %d\n", id)
	return u.CodeSnippetStorage.GetCodeSnippetById(id)
}

func validateLanguage(lang string) error {
	// TODO
	return nil
}
