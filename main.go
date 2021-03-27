package main

import (
	"flag"
	"fmt"
	accountstorage "github.com/mp-hl-2021/code-swamp/accountStorage"
	"github.com/mp-hl-2021/code-swamp/api"
	"github.com/mp-hl-2021/code-swamp/auth"
	"github.com/mp-hl-2021/code-swamp/usecases"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	privateKeyPath := flag.String("privateKey", "app.rsa", "file path")
	publicKeyPath := flag.String("publicKey", "app.rsa.pub", "file path")
	flag.Parse()

	privateKeyBytes, err := ioutil.ReadFile(*privateKeyPath)
	publicKeyBytes, err := ioutil.ReadFile(*publicKeyPath)

	a, err := auth.NewJwt(privateKeyBytes, publicKeyBytes, 100*time.Minute)
	if err != nil {
		panic(err)
	}

	user := &usecases.User{
		AccountStorage: accountstorage.NewMemory(),
		Auth: a,
	}

	service := api.NewApi(user)

	addr := "localhost:8080"
	server := http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: service.Router(),
	}
	fmt.Printf("Serving at %s\n", addr)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}