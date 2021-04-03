package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/mp-hl-2021/code-swamp/internal/usecases"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type AccountFake struct{}

func (AccountFake) CreateAccount(login, password string) (usecases.Account, error) {
	if login == "katyukha" {
		return usecases.Account{Id: 1}, nil
	}
	if password == "  " {
		return usecases.Account{}, usecases.ErrInvalidPasswordString
	}
	return usecases.Account{}, errors.New("failed to create account")
}

func (AccountFake) LoginToAccount(login, password string) (string, error) {
	if login == "katyukha" && password == "kek1234" {
		return "token", nil
	}
	if login == "masha" && password != "123" {
		return "", usecases.ErrInvalidPassword
	}
	if password == "  " {
		return "", usecases.ErrInvalidPasswordString
	}
	return "", errors.New("failed to login to account")
}

func (AccountFake) GetMyLinks(a usecases.Account) ([]string, error) {
	panic("not implemented")
}

func (AccountFake) CreateSnippet(a *usecases.Account, code string, lang *string, lifetime time.Duration) (string, error) {
	panic("not implemented")
}

func (AccountFake) GetSnippetById(string) (usecases.CodeSnippet, error) {
	panic("not implemented")
}

func (AccountFake) GetAccountByToken(string) (usecases.Account, error) {
	panic("not implemented")
}

func assertStatusCode(t *testing.T, expectedCode, actualCode int) {
	if expectedCode != actualCode {
		t.Errorf("Server MUST return %d (%s) status code, but %d (%s) given",
			expectedCode, http.StatusText(expectedCode), actualCode, http.StatusText(actualCode))
	}
}

func invalidJsonTest(router http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte("{a:")))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func makeSignupRequest(t *testing.T, router http.Handler, login, password string) *httptest.ResponseRecorder {
	m := postSignupRequestModel{
		Login:    login,
		Password: password,
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal("failed to marshal struct")
	}
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(b))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	return resp
}

func makeSigninRequest(t *testing.T, router http.Handler, login, password string) *httptest.ResponseRecorder {
	m := postSignupRequestModel{
		Login:    login,
		Password: password,
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal("failed to marshal struct")
	}
	req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewReader(b))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	return resp
}

func Test_postSignup(t *testing.T) {
	service := NewApi(&AccountFake{})
	router := service.Router()

	t.Run("failure on invalid json", func(t *testing.T) {
		resp := invalidJsonTest(router, "/signup")
		assertStatusCode(t, resp.Code, http.StatusBadRequest)
	})

	t.Run("failed to create account", func(t *testing.T) {
		resp := makeSignupRequest(t, router, "masha", "123")
		assertStatusCode(t, resp.Code, http.StatusInternalServerError)
	})

	t.Run("successful account creation", func(t *testing.T) {
		resp := makeSignupRequest(t, router, "katyukha", "kek1234")
		assertStatusCode(t, resp.Code, http.StatusCreated)
	})

	t.Run("failure on invalid password string", func(t *testing.T) {
		resp := makeSignupRequest(t, router, "jaba", "  ")
		assertStatusCode(t, resp.Code, http.StatusBadRequest)
	})
}

func Test_postSignin(t *testing.T) {
	service := NewApi(&AccountFake{})
	router := service.Router()

	t.Run("failure on invalid json", func(t *testing.T) {
		resp := invalidJsonTest(router, "/signup")
		assertStatusCode(t, resp.Code, http.StatusBadRequest)
	})
	t.Run("failed login with incorrect login or password", func(t *testing.T) {
		resp := makeSigninRequest(t, router, "masha", "0")
		assertStatusCode(t, resp.Code, http.StatusUnauthorized)
	})
	t.Run("successful login with correct password", func(t *testing.T) {
		resp := makeSigninRequest(t, router, "katyukha", "kek1234")
		assertStatusCode(t, resp.Code, http.StatusOK)
	})
	t.Run("failure on invalid password string", func(t *testing.T) {
		resp := makeSigninRequest(t, router, "jaba", "  ")
		assertStatusCode(t, resp.Code, http.StatusBadRequest)
	})
}
