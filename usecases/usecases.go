package usecases

import (
	"fmt"
	"time"
)

type Account struct {
	Id string
}

type CodeSnippet struct {
	Code     string
	Lang     *string
	Lifetime time.Duration
}

type UserInterface interface {
	CreateAccount(login, password string) (Account, error)
	LoginToAccount(login, password string) (string, error)

	GetMyLinks(a Account) ([]string, error)
	CreateSnippet(a *Account, cs CodeSnippet) (string, error)
	GetSnippetById(string) (CodeSnippet, error)
}

type User struct{}

func (User) CreateAccount(login, password string) (Account, error) {
	// TODO
	fmt.Printf("Register: %s %s", login, password)
	return Account{Id: "0"}, nil
}

func (User) LoginToAccount(login, password string) (string, error) {
	// TODO
	fmt.Printf("Login: %s %s", login, password)
	return "token", nil
}

func (User) GetMyLinks(a Account) ([]string, error) {
	// TODO
	fmt.Printf("GetMyLinks: %s", a.Id)
	return []string{"a", "b", "c"}, nil
}

func (User) CreateSnippet(a *Account, cs CodeSnippet) (string, error) {
	// TODO
	fmt.Printf("CrateLink: %s %s", a.Id, cs.Code)
	return "id", nil
}

func (User) GetSnippetById(id string) (CodeSnippet, error) {
	// TODO
	fmt.Printf("GetSnippetById: %s", id)
	return CodeSnippet{Code: "code"}, nil
}
