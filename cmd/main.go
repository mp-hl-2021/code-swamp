package main

import (
	"flag"
	"fmt"
	"github.com/mp-hl-2021/code-swamp/internal/interface/httpapi"
	"github.com/mp-hl-2021/code-swamp/internal/interface/memory/accountrepo"
	"github.com/mp-hl-2021/code-swamp/internal/interface/memory/codesnippetrepo"
	"github.com/mp-hl-2021/code-swamp/internal/service/token"
	"github.com/mp-hl-2021/code-swamp/internal/usecases/account"
	"github.com/mp-hl-2021/code-swamp/internal/usecases/codesnippet"
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

	a, err := token.NewJwt(privateKeyBytes, publicKeyBytes, 100*time.Minute)
	if err != nil {
		panic(err)
	}

	accountUseCases := &account.UseCases{
		AccountStorage: accountrepo.NewMemory(),
		Auth:           a,
	}

	codeSnippetUseCases := &codesnippet.UseCases{
		CodeSnippetStorage: codesnippetrepo.NewMemory(),
	}

	service := httpapi.NewApi(accountUseCases, codeSnippetUseCases)

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