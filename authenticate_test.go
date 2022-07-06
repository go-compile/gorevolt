package gorevolt_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-compile/gorevolt"
)

func TestAuthenticate(t *testing.T) {
	token, exists := os.LookupEnv("gorevolt_test_token")
	if !exists {
		t.Skip("no token provided in `gorevolt_test_token` environment variable")
	}

	client := gorevolt.New(token)

	client.Register(func(c *gorevolt.Client, startup time.Duration) {
		fmt.Printf("[CONNECTED] [USER: %s]\n", c.User.Username)
	})

	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Minute)

	client.Close()
}
