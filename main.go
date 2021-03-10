package main

import (
	"fmt"
	"github.com/mp-hl-2021/code-swamp/api"
	"github.com/mp-hl-2021/code-swamp/usecases"
	"net/http"
	"time"
)

func main() {
	user := &usecases.User{}

	service := api.NewApi(user)

	addr := "localhost:8080"
	server := http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	fmt.Printf("Serving at %s\n", addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}