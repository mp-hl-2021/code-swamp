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
	if a.Id == 1 {
		return []string{"a", "b", "c"}, nil
	}
	return []string{}, errors.New("failed to get links")
}

func (AccountFake) CreateSnippet(a *usecases.Account, code string, lang string, lifetime time.Duration) (uint, error) {
	if lang == "petooh" {
		return 0, usecases.ErrInvalidLanguage
	}
	if code == "internal" {
		return 0, errors.New("failed ti create new snippet")
	}
	return 1, nil
}

func (AccountFake) GetSnippetById(string) (usecases.CodeSnippet, error) {
	panic("not implemented")
}

func (AccountFake) GetAccountByToken(token string) (usecases.Account, error) {
	if token == "correct" {
		return usecases.Account{Id: 1}, nil
	}
	if token == "internal" {
		return usecases.Account{Id: 100}, nil
	}
	return usecases.Account{}, errors.New("invalid token claims")
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

func Test_postSignup(t *testing.T) {
	service := NewApi(&AccountFake{})
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
	service := NewApi(&AccountFake{})
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
	service := NewApi(&AccountFake{})
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
	service := NewApi(&AccountFake{})
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