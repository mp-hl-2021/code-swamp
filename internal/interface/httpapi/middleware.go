package httpapi

import (
	"context"
	"net/http"
	"strings"
)

func (a *Api) authenticate(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearHeader := r.Header.Get("Authorization")
		strArr := strings.Split(bearHeader, " ")
		if len(strArr) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token := strArr[1]
		id, err := a.AccountUseCases.Authenticate(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), accountIdContextKey, id)
		handler(w, r.WithContext(ctx))
	}
}


func (a *Api) authenticateOrNot(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearHeader := r.Header.Get("Authorization")
		strArr := strings.Split(bearHeader, " ")
		if len(strArr) != 2 {
			handler(w, r)
		} else {
			token := strArr[1]
			id, err := a.AccountUseCases.Authenticate(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), accountIdContextKey, id)
			handler(w, r.WithContext(ctx))
		}
	}
}
