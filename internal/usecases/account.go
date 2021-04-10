package usecases

import (
	"errors"
	"fmt"
	account "github.com/mp-hl-2021/code-swamp/internal/domain/account"
	"github.com/mp-hl-2021/code-swamp/internal/domain/codesnippet"
	"unicode"

	"golang.org/x/crypto/bcrypt"

	"time"
)

var (
	ErrInvalidLoginString    = errors.New("login string contains invalid character")
	ErrInvalidLoginString2   = errors.New("login string should start with a letter")
	ErrInvalidPasswordString = errors.New("password string contains invalid character")
	ErrTooShortString        = errors.New("too short string")
	ErrTooLongString         = errors.New("too long string")
	ErrNoDigits              = errors.New("password string contains no digits")
	ErrNoUpperCaseLetters    = errors.New("password string contains no upper case letters")
	ErrNoLowerCaseLetters    = errors.New("password string contains no lower case letters")

	ErrInvalidLogin    = errors.New("login not found")
	ErrInvalidPassword = errors.New("invalid password")

	ErrInvalidLanguage = errors.New("language is invalid")
)

const (
	minLoginLength    = 6
	maxLoginLength    = 30
	minPasswordLength = 6
	maxPasswordLength = 40
)

type Account struct {
	Id uint
}

type AccountInterface interface {
	CreateAccount(login, password string) (Account, error)
	LoginToAccount(login, password string) (string, error)

	GetMySnippetIds(a Account) ([]uint, error)
	CreateSnippet(a *Account, code string, lang string, lifetime time.Duration) (uint, error)
	GetSnippetById(uint) (codesnippet.CodeSnippet, error)
	GetAccountByToken(string) (Account, error)
}

type User struct {
	Auth           Interface
	AccountStorage account.Interface
	CodeSnippetStorage codesnippet.Interface
}

func (u *User) CreateAccount(login, password string) (Account, error) {
	fmt.Printf("Register: %s %s\n", login, password)
	if err := validateLogin(login); err != nil {
		return Account{}, err
	}
	if err := validatePassword(password); err != nil {
		return Account{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, err
	}

	acc, err := u.AccountStorage.CreateAccount(account.Credentials{
		Login:    login,
		Password: string(hashedPassword),
	})
	if err != nil {
		return Account{}, err
	}
	return Account{Id: acc.Id}, nil
}

func (u *User) LoginToAccount(login, password string) (string, error) {
	fmt.Printf("Login: %s %s\n", login, password)
	if err := validateLogin(login); err != nil {
		return "", err
	}
	if err := validatePassword(password); err != nil {
		return "", err
	}
	acc, err := u.AccountStorage.GetAccountByLogin(login)
	if err != nil {
		return "", ErrInvalidLogin
	}
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Credentials.Password), []byte(password)); err != nil {
		return "", ErrInvalidPassword
	}

	token, err := u.Auth.IssueToken(acc.Id)

	return token, err
}

func (u *User) GetMySnippetIds(a Account) ([]uint, error) {
	fmt.Printf("GetMySnippetIds: %s", a.Id)
	return u.CodeSnippetStorage.GetMyCodeSnippetIds(a.Id)
}

func (u *User) CreateSnippet(a *Account, code string, lang string, lifetime time.Duration) (uint, error) {
	fmt.Printf("CreateSnippet: %s %s", a.Id, code)
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

func (u *User) GetSnippetById(id uint) (codesnippet.CodeSnippet, error) {
	fmt.Printf("GetSnippetById: %s", id)
	return u.CodeSnippetStorage.GetCodeSnippetById(id)
}

func (u User) GetAccountByToken(token string) (Account, error) {
	id, err := u.Auth.UserIdByToken(token)
	if err != nil {
		return Account{}, err
	}
	wAcc, err := u.AccountStorage.GetAccountById(id)
	return Account{wAcc.Id}, err
}

func validateLogin(login string) error {
	chars := 0
	if !unicode.IsLetter([]rune(login)[0]) {
		return ErrInvalidLoginString2
	}
	for _, r := range login {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ErrInvalidLoginString
		}
		chars++
	}
	if chars < minLoginLength {
		return ErrTooShortString
	}
	if chars > maxLoginLength {
		return ErrTooLongString
	}
	return nil
}

func validatePassword(password string) error {
	chars := 0
	lower := 0
	upper := 0
	digit := 0
	for _, r := range password {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return ErrInvalidPasswordString
		}
		if unicode.IsLower(r) {
			lower = 1
		}

		if unicode.IsUpper(r) {
			upper = 1
		}

		if unicode.IsDigit(r) {
			digit = 1
		}

		chars++
	}
	if chars < minPasswordLength {
		return ErrTooShortString
	}
	if chars > maxPasswordLength {
		return ErrTooLongString
	}
	if lower == 0 {
		return ErrNoLowerCaseLetters
	}
	if upper == 0 {
		return ErrNoUpperCaseLetters
	}
	if digit == 0 {
		return ErrNoDigits
	}
	return nil
}

func validateLanguage(lang string) error {
	// TODO
	return nil
}
