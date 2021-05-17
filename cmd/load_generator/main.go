package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mp-hl-2021/code-swamp/internal/interface/httpapi"
	"github.com/mp-hl-2021/code-swamp/internal/usecases/codesnippet"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var config = struct {
	address          string
	concurrencyLevel int
}{}

func init() {
	address := flag.String("address", "http://localhost:8080", "swamp address")
	concurrencyLevel := flag.Int("concurrency", 50, "a number of concurrent requests")
	flag.Parse()

	config.address = *address
	config.concurrencyLevel = *concurrencyLevel
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	defer func() {
		signal.Stop(ch)
		cancel()
	}()

	go func() {
		select {
		case <-ch:
			cancel()
		case <-ctx.Done():
		}
	}()

	c := client{
		c: http.Client{
			Timeout: 10 * time.Second,
		},
	}

	var wg sync.WaitGroup
	wg.Add(config.concurrencyLevel)
	for i := 0; i < config.concurrencyLevel; i++ {
		go func(i int) {
			var err error
			if i % 2 == 0 {
				err = accountCreator(ctx, c)
			} else {
				err = snippetCreator(ctx, c)
			}
			fmt.Printf("worker %d finished, err: %v\n", i, err)
			wg.Done()

		}(i)
	}
	wg.Wait()
	fmt.Println("all workers have been finished")
}

func accountCreator(ctx context.Context, c client) error {
	for {
		select {
		default:
			login := gofakeit.Username() + gofakeit.DigitN(9)
			password := gofakeit.Password(true, true, true, false, false, 16)
			err := c.createAccount(ctx, login, password)
			if err != nil {
				fmt.Println("request failed:", err)
			}
		case <-ctx.Done():
			fmt.Println("leaving worker")
			return ctx.Err()
		}
	}
}

func snippetCreator(ctx context.Context, c client) error {
	for {
		select {
		default:
			code := gofakeit.LetterN(20)
			lang := gofakeit.RandomString(codesnippet.SupportedLanguages)
			_, err := c.createCodeSnippet(ctx, code, lang)
			if err != nil {
				fmt.Println("request failed:", err)
			}
		case <-ctx.Done():
			fmt.Println("leaving worker")
			return ctx.Err()
		}
	}
}

type client struct {
	c http.Client
}

func (c client) createAccount(ctx context.Context, login, password string) error {
	body := httpapi.PostSignupRequestModel{Login: login, Password: password}
	s, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.address+"/signup", bytes.NewReader(s))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create account: %v", resp.Status)
	}
	return nil
}

func (c client) createCodeSnippet(ctx context.Context, code, lang string) (string, error) {
	body := httpapi.PostCodeRequestModel  {Code: code, Lang: lang, Lifetime: time.Millisecond}
	s, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.address+"/", bytes.NewReader(s))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed to create code snippet: %v", resp.Status)
	}
	return resp.Header.Get("Location"), nil
}