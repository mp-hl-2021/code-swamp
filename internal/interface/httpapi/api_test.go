package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mp-hl-2021/code-swamp/internal/domain/codesnippet"
	"github.com/mp-hl-2021/code-swamp/internal/interface/memory/codesnippetrepo"
	"github.com/mp-hl-2021/code-swamp/internal/usecases/account"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type AccountFake struct{}
type CodeSnippetFake struct{}

func (AccountFake) CreateAccount(login, password string) (account.Account, error) {
	if login == "katyukha" {
		return account.Account{Id: 1}, nil
	}
	if password == "  " {
		return account.Account{}, account.ErrInvalidPasswordString
	}
	return account.Account{}, errors.New("failed to create account")
}

func (AccountFake) LoginToAccount(login, password string) (string, error) {
	if login == "katyukha" && password == "kek1234" {
		return "token", nil
	}
	if login == "masha" && password != "123" {
		return "", account.ErrInvalidPassword
	}
	if password == "  " {
		return "", account.ErrInvalidPasswordString
	}
	return "", errors.New("failed to login to account")
}

func (CodeSnippetFake) GetMySnippetIds(a account.Account) ([]uint, error) {
	if a.Id == 1 {
		return []uint{1, 2, 3}, nil
	}
	return []uint{}, errors.New("failed to get links")
}

func (CodeSnippetFake) CreateSnippet(a *account.Account, code string, lang string, lifetime time.Duration) (uint, error) {
	if lang == "petooh" {
		return 0, account.ErrInvalidLanguage
	}
	if code == "internal" {
		return 0, errors.New("failed ti create new snippet")
	}
	return 1, nil
}

func (CodeSnippetFake) GetSnippetById(sid uint) (codesnippet.CodeSnippet, error) {
	if sid == 1 {
		return codesnippet.CodeSnippet{}, codesnippetrepo.ErrInvalidSnippedId
	}
	if sid == 2 {
		return codesnippet.CodeSnippet{}, errors.New("failed to get snippet")
	}
	return codesnippet.CodeSnippet{Code: "KoKoKoKoKoKoKoKoKoKo Kud-Kudah"}, nil
}

func (AccountFake) GetAccountByToken(token string) (account.Account, error) {
	if token == "correct" {
		return account.Account{Id: 1}, nil
	}
	if token == "internal" {
		return account.Account{Id: 100}, nil
	}
	return account.Account{}, errors.New("invalid token claims")
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

func makeGetLinksRequest(t *testing.T, router http.Handler, token string) *httptest.ResponseRecorder {
	m := postLinksRequestModel{
		Token: token,
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal("failed to marshal struct")
	}
	req := httptest.NewRequest(http.MethodPost, "/myswamp", bytes.NewReader(b))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	return resp
}

func makePostCodeRequest(t *testing.T, router http.Handler, token, code, lang string) *httptest.ResponseRecorder {
	d, _ := time.ParseDuration("12h")
	m := postCodeRequestModel{
		Token:    token,
		Code:     code,
		Lang:     lang,
		Lifetime: d,
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal("failed to marshal struct")
	}
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	return resp
}

func makeGetCodeRequest(router http.Handler, sid uint) *httptest.ResponseRecorder  {
	req := httptest.NewRequest(http.MethodGet, "/toad/"+fmt.Sprintf("%d",sid), bytes.NewReader([]byte("")))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func Test_postSignup(t *testing.T) {
	service := NewApi(&AccountFake{}, &CodeSnippetFake{})
	router := service.Router()

	t.Run("failure on invalid json", func(t *testing.T) {
		resp := invalidJsonTest(router, "/signup")
		assertStatusCode(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("failed to create account", func(t *testing.T) {
		resp := makeSignupRequest(t, router, "masha", "123")
		assertStatusCode(t, http.StatusInternalServerError, resp.Code)
	})

	t.Run("successful account creation", func(t *testing.T) {
		resp := makeSignupRequest(t, router, "katyukha", "kek1234")
		assertStatusCode(t, http.StatusCreated, resp.Code)
	})

	t.Run("failure on invalid password string", func(t *testing.T) {
		resp := makeSignupRequest(t, router, "jaba", "  ")
		assertStatusCode(t, http.StatusBadRequest, resp.Code)
	})
}

func Test_postSignin(t *testing.T) {
	service := NewApi(&AccountFake{}, &CodeSnippetFake{})
	router := service.Router()

	t.Run("failure on invalid json", func(t *testing.T) {
		resp := invalidJsonTest(router, "/signup")
		assertStatusCode(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("failed login with incorrect login or password", func(t *testing.T) {
		resp := makeSigninRequest(t, router, "masha", "0")
		assertStatusCode(t, http.StatusUnauthorized, resp.Code)
	})
	t.Run("successful login with correct password", func(t *testing.T) {
		resp := makeSigninRequest(t, router, "katyukha", "kek1234")
		assertStatusCode(t, http.StatusOK, resp.Code)
	})
	t.Run("failure on invalid password string", func(t *testing.T) {
		resp := makeSigninRequest(t, router, "jaba", "  ")
		assertStatusCode(t, http.StatusBadRequest, resp.Code)
	})
}

func Test_postLinks(t *testing.T) {
	service := NewApi(&AccountFake{}, &CodeSnippetFake{})
	router := service.Router()

	t.Run("failure on invalid json", func(t *testing.T) {
		resp := invalidJsonTest(router, "/myswamp")
		assertStatusCode(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("failed to get links with incorrect token", func(t *testing.T) {
		resp := makeGetLinksRequest(t, router, "incorrect")
		assertStatusCode(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("failed to get links for existing user", func(t *testing.T) {
		resp := makeGetLinksRequest(t, router, "internal")
		assertStatusCode(t, http.StatusInternalServerError, resp.Code)
	})
	t.Run("successful obtainment of links ", func(t *testing.T) {
		resp := makeGetLinksRequest(t, router, "correct")
		assertStatusCode(t, http.StatusOK, resp.Code)
	})
}

func Test_postCode(t *testing.T) {
	service := NewApi(&AccountFake{}, &CodeSnippetFake{})
	router := service.Router()

	t.Run("failure on invalid json", func(t *testing.T) {
		resp := invalidJsonTest(router, "/")
		assertStatusCode(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("failed to post code with invalid token", func(t *testing.T) {
		resp := makePostCodeRequest(t, router, "incorrect", "",  "")
		assertStatusCode(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("failed to post code with invalid language", func(t *testing.T) {
		resp := makePostCodeRequest(t, router, "correct", "KoKoKoKoKoKoKoKoKoKo Kud-Kudah", "petooh")
		assertStatusCode(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("failed to post code because of internal error", func(t *testing.T) {
		resp := makePostCodeRequest(t, router, "", "internal", "")
		assertStatusCode(t, http.StatusInternalServerError, resp.Code)
	})
	t.Run("successful snippet creation ", func(t *testing.T) {
		resp := makePostCodeRequest(t, router, "", "KoKoKoKoKoKoKoKoKoKo Kud-Kudah", "")
		assertStatusCode(t, http.StatusCreated, resp.Code)
	})
}

func Test_getCode(t *testing.T) {
	service := NewApi(&AccountFake{}, &CodeSnippetFake{})
	router := service.Router()

	t.Run("failure on invalid json", func(t *testing.T) {
		resp := invalidJsonTest(router, "/")
		assertStatusCode(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("no such code snipped", func(t *testing.T) {
		resp := makeGetCodeRequest(router, 1)
		assertStatusCode(t, http.StatusBadRequest, resp.Code)
	})
	t.Run("failed to get snippet", func(t *testing.T) {
		resp := makeGetCodeRequest(router, 2)
		assertStatusCode(t, http.StatusInternalServerError, resp.Code)
	})
	t.Run("successful obtainment of snippet", func(t *testing.T) {
		resp := makeGetCodeRequest(router, 3)
		assertStatusCode(t, http.StatusOK, resp.Code)
	})
}