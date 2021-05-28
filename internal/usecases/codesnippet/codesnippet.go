package codesnippet

import (
	"errors"
	"fmt"
	"github.com/mp-hl-2021/code-swamp/internal/domain/codesnippet"
	"github.com/mp-hl-2021/code-swamp/internal/usecases/account"
	"io/ioutil"
	"os"
	"os/exec"
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

type CheckCodeRequest struct {
	Sid  uint
	Code string
	Lang string
}

type Interface interface {
	GetMySnippetIds(a account.Account) ([]uint, error)
	CreateSnippet(a *account.Account, code string, lang string, lifetime time.Duration) (uint, error)
	GetSnippetById(uint) (codesnippet.CodeSnippet, error)
	CheckCode(sid uint, code string, lang string) error
}

type UseCases struct {
	CodeSnippetStorage codesnippet.Interface
	CodeCheckChannel   chan<- CheckCodeRequest
}

func RunLinter(code string, lang string) (CodeCheckResult, error) {
	file, err := ioutil.TempFile("", "tmp")
	if err != nil {
		return CodeCheckResult{}, errors.New("failed to create temporary file")
	}
	defer os.Remove(file.Name())
	_, err = file.Write([]byte(code))
	if err != nil {
		return CodeCheckResult{}, errors.New("failed to write to temporary file")
	}
	output, err := exec.Command("dupl", "-t", "100", file.Name()).Output()
	if err != nil {
		return CodeCheckResult{}, errors.New("failed to run dupl on file")
	}
	return CodeCheckResult{correct: true, msg: string(output)}, nil
}

func (u *UseCases) CheckCode(sid uint, code string, lang string) error {
	r, err := RunLinter(code, lang)
	var status bool
	var msg string
	if err != nil {
		status = false
		msg = err.Error()
	} else if !r.correct {
		status = false
		msg = r.msg
	} else {
		status = true
		msg = ""
	}
	return u.CodeSnippetStorage.SetCodeStatus(sid, status, msg)
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
		u.CodeCheckChannel <- CheckCodeRequest{sid, code, lang}
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
