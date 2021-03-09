package main

import (
	"github.com/mp-hl-2021/code-swamp/api"
	"github.com/mp-hl-2021/code-swamp/usecases"
	"net/http"
	"time"
)

func main() {
	user := &usecases.User{}

	service := api.NewApi(user)

	server := http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}