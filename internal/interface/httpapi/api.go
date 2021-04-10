package httpapi

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mp-hl-2021/code-swamp/internal/interface/memory/codesnippetrepo"
	"github.com/mp-hl-2021/code-swamp/internal/usecases/account"
	"github.com/mp-hl-2021/code-swamp/internal/usecases/codesnippet"
	"net/http"
	"strconv"
	"time"
)

const snippetIdUrlPathKey = "snippet_id"

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

	router.HandleFunc("/signup", a.postSignup).Methods(http.MethodPost)
	router.HandleFunc("/signin", a.postSignin).Methods(http.MethodPost)

	router.HandleFunc("/myswamp", a.postLinks).Methods(http.MethodPost)
	router.HandleFunc("/", a.postCode).Methods(http.MethodPost)

	router.HandleFunc("/toad/{"+snippetIdUrlPathKey+"}", a.getCode).Methods(http.MethodGet)
	router.HandleFunc("/toad/{"+snippetIdUrlPathKey+"}/download", a.getCodeFile).Methods(http.MethodGet)

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
			account.ErrInvalidLoginString,
			account.ErrInvalidLoginString2,
			account.ErrInvalidPasswordString,
			account.ErrTooShortString,
			account.ErrTooLongString,
			account.ErrNoDigits,
			account.ErrNoUpperCaseLetters,
			account.ErrNoLowerCaseLetters:

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
	w.Write([]byte(token))
}

type postLinksRequestModel struct {
	Token string `json:"token"`
}

type postLinksResponseModel struct {
	Links []string `json:"links"`
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

	ss, err := a.CodeSnippetUseCases.GetMySnippetIds(acc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mm := postLinksResponseModel{
		Links: make([]string, len(ss)),
	}
	for i := range ss {
		mm.Links[i] = fmt.Sprintf("/%d", ss[i])
	}
	if err := json.NewEncoder(w).Encode(mm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

	var acc *account.Account = nil
	if m.Token != "" {
		a, err := a.AccountUseCases.GetAccountByToken(m.Token)
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
		return
	}

	location := fmt.Sprintf("/%d", id)
	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)
}

type getCodeResponseModel struct {
	Code string `json:"code"`
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
	m := getCodeResponseModel{
		Code: ss.Code,
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Api) getCodeFile(w http.ResponseWriter, _ *http.Request) {
	// TODO: generate snippet id by url.
	w.WriteHeader(http.StatusOK)
}
