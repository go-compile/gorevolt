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
	close := make(chan struct{})

	client.Register(func(c *gorevolt.Client, startup time.Duration) {
		fmt.Printf("[CONNECTED] [USER: %s]\n", c.User.Username)

		// close and finish test
		close <- struct{}{}
	})

	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}

	<-close
	client.Close()
}
