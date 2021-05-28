package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	repository "github.com/mp-hl-2021/code-swamp/internal/domain/account"
	"github.com/mp-hl-2021/code-swamp/internal/interface/memory/codesnippetrepo"
	"github.com/mp-hl-2021/code-swamp/internal/interface/prom"
	"github.com/mp-hl-2021/code-swamp/internal/usecases/account"
	"github.com/mp-hl-2021/code-swamp/internal/usecases/codesnippet"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

const (
	accountIdContextKey = "account_id"
	snippetIdUrlPathKey = "snippet_id"
)

type Api struct {
	AccountUseCases     account.Interface
	CodeSnippetUseCases codesnippet.Interface
}

func NewApi(a account.Interface, c codesnippet.Interface) *Api {
	return &Api{
		AccountUseCases:     a,
		CodeSnippetUseCases: c,
	}
}

func (a *Api) Router() http.Handler {
	router := mux.NewRouter()

	router.Use(prom.Measurer())
	router.Use(a.logger)

	router.HandleFunc("/signup", a.postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", a.postSignin).Methods(http.MethodPost)

	router.HandleFunc("/myswamp", a.authenticate(a.postLinks)).Methods(http.MethodPost)
	router.HandleFunc("/", a.authenticateOrNot(a.postCode)).Methods(http.MethodPost)

	router.HandleFunc("/toad/{"+snippetIdUrlPathKey+"}", a.getCode).Methods(http.MethodGet)

	router.Handle("/metrics", promhttp.Handler())

	return router
}

type PostSignupRequestModel struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (a *Api) postSignup(w http.ResponseWriter, r *http.Request) {
	var m PostSignupRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := a.AccountUseCases.CreateAccount(m.Login, m.Password)
	if err != nil {
		var statusCode int
		switch err {
		case
			account.ErrInvalidLoginString,
			account.ErrInvalidLoginString2,
			account.ErrInvalidPasswordString,
			account.ErrTooShortString,
			account.ErrTooLongString,
			account.ErrNoDigits,
			account.ErrNoUpperCaseLetters,
			account.ErrNoLowerCaseLetters,
			repository.ErrAlreadyExist:

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
	var m PostSignupRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := a.AccountUseCases.LoginToAccount(m.Login, m.Password)
	if err != nil {
		var statusCode int
		switch err {

		case
			account.ErrInvalidLogin,
			account.ErrInvalidPassword:

			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusBadRequest
		}
		w.WriteHeader(statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/jwt")
	if _, err = w.Write([]byte(token)); err != nil {
		return
	}
}

type PostLinksResponseModel struct {
	Links []string `json:"links"`
}

func (a *Api) postLinks(w http.ResponseWriter, r *http.Request) {
	aid, ok := r.Context().Value(accountIdContextKey).(uint)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	acc, err := a.AccountUseCases.GetAccountById(aid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ss, err := a.CodeSnippetUseCases.GetMySnippetIds(acc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mm := PostLinksResponseModel{
		Links: make([]string, len(ss)),
	}
	for i := range ss {
		mm.Links[i] = fmt.Sprintf("/toad/%d", ss[i])
	}
	if err := json.NewEncoder(w).Encode(mm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type PostCodeRequestModel struct {
	Code     string        `json:"code"`
	Lang     string        `json:"lang"`
	Lifetime time.Duration `json:"lifetime"`
}

func (a *Api) postCode(w http.ResponseWriter, r *http.Request) {
	var m PostCodeRequestModel
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var acc *account.Account = nil
	aid, ok := r.Context().Value(accountIdContextKey).(uint)
	if ok {
		a, err := a.AccountUseCases.GetAccountById(aid)
		acc = &a
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	id, err := a.CodeSnippetUseCases.CreateSnippet(acc, m.Code, m.Lang, m.Lifetime)
	if err != nil {
		var statusCode int
		switch err {

		case
			account.ErrInvalidLanguage:

			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		fmt.Println(err)
		return
	}

	location := fmt.Sprintf("/toad/%d", id)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

func (a *Api) getCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s, ok := vars[snippetIdUrlPathKey]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sid, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ss, err := a.CodeSnippetUseCases.GetSnippetById(uint(sid))
	if err != nil {
		var statusCode int
		switch err {

		case
			codesnippetrepo.ErrInvalidSnippedId:

			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		w.WriteHeader(statusCode)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	status := ""
	if !ss.IsChecked {
		status = "Not checked yet"
	} else if ss.IsCorrect {
		status = "Correct"
	} else {
		status = "Incorrect"
	}
	w.Write([]byte("Status:" + status + ", Code: " + ss.Code))
}
