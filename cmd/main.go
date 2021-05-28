package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mp-hl-2021/code-swamp/internal/interface/httpapi"
	"github.com/mp-hl-2021/code-swamp/internal/interface/postgres/accountrepo"
	"github.com/mp-hl-2021/code-swamp/internal/interface/postgres/codesnippetrepo"
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

	// TODO: pass arguments with config
	connStr := "user=postgres password=12345 host=db dbname=postgres sslmode=disable"

	conn, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	ch := make(chan codesnippet.CheckCodeRequest)

	accountUseCases := &account.UseCases{
		AccountStorage: accountrepo.New(conn),
		Auth:           a,
	}

	codeSnippetUseCases := &codesnippet.UseCases{
		CodeSnippetStorage: codesnippetrepo.New(conn),
		CodeCheckChannel:   ch,
	}

	go func() {
		for _ = range time.Tick(time.Minute) {
			c := <-ch
			fmt.Printf("Checking code sid: %d, code: %s, lang: %s\n", c.Sid, c.Code, c.Lang)
			err := codeSnippetUseCases.CheckCode(c.Sid, c.Code, c.Lang)
			if err != nil {
				fmt.Printf("Error checking code: %s\n", err)
			}
		}
	}()

	service := httpapi.NewApi(accountUseCases, codeSnippetUseCases)

	addr := ":8080"
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