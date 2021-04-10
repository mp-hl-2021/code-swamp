package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/code-swamp/internal/usecases"
	"net/http"
	"strings"
	"time"
)

type Api struct {
	AccountUseCases usecases.AccountInterface
}

func NewApi(u usecases.AccountInterface) *Api {
	return &Api{
		AccountUseCases: u,
	}
}

func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/signup", a.postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", a.postSignin).Methods(http.MethodPost)

	router.HandleFunc("/myswamp", a.postLinks).Methods(http.MethodPost)
	router.HandleFunc("/", a.postCode).Methods(http.MethodPost)

	router.HandleFunc("/{snippet_id}", a.getCode).Methods(http.MethodGet)
	router.HandleFunc("/{snippet_id}/download", a.getCodeFile).Methods(http.MethodGet)

	return router
}

type postSignupRequestModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (a *Api) postSignup(w http.ResponseWriter, r *http.Request) {
	var m postSignupRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := a.AccountUseCases.CreateAccount(m.Login, m.Password)
	if err != nil {
		var statusCode int
		switch err {
		case
			usecases.ErrInvalidLoginString,
			usecases.ErrInvalidLoginString2,
			usecases.ErrInvalidPasswordString,
			usecases.ErrTooShortString,
			usecases.ErrTooLongString,
			usecases.ErrNoDigits,
			usecases.ErrNoUpperCaseLetters,
			usecases.ErrNoLowerCaseLetters:

			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *Api) postSignin(w http.ResponseWriter, r *http.Request) {
	var m postSignupRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := a.AccountUseCases.LoginToAccount(m.Login, m.Password)
	if err != nil {
		var statusCode int
		switch err {

		case
			usecases.ErrInvalidLogin,
			usecases.ErrInvalidPassword:

			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusBadRequest
		}
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/jwt")
	w.Write([]byte(token))
}

type postLinksRequestModel struct {
	Token string `json:"token"`
}

func (a *Api) postLinks(w http.ResponseWriter, r *http.Request) {
	var m postLinksRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	acc, err := a.AccountUseCases.GetAccountByToken(m.Token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	links, err := a.AccountUseCases.GetMyLinks(acc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(strings.Join(links, ",")))
}

type postCodeRequestModel struct {
	Token    string        `json:"token"`
	Code     string        `json:"code"`
	Lang     string        `json:"lang"`
	Lifetime time.Duration `json:"lifetime"`
}

func (a *Api) postCode(w http.ResponseWriter, r *http.Request) {
	var m postCodeRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var acc *usecases.Account = nil
	if m.Token != "" {
		a, err := a.AccountUseCases.GetAccountByToken(m.Token)
		acc = &a
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	id, err := a.AccountUseCases.CreateSnippet(acc, m.Code, m.Lang, m.Lifetime)
	if err != nil {
		var statusCode int
		switch err {

		case
			usecases.ErrInvalidLanguage:

			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		return
	}

	location := fmt.Sprintf("/%d", id)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

func (a *Api) getCode(w http.ResponseWriter, _ *http.Request) {
	// TODO: generate snippet id by url.
	w.WriteHeader(http.StatusOK)
}

func (a *Api) getCodeFile(w http.ResponseWriter, _ *http.Request) {
	// TODO: generate snippet id by url.
	w.WriteHeader(http.StatusOK)
}
