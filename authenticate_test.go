package gorevolt_test

import (
	"os"
	"testing"

	"github.com/go-compile/gorevolt"
)

func TestAuthenticate(t *testing.T) {
	token, exists := os.LookupEnv("gorevolt_test_token")
	if !exists {
		t.Skip("no token provided in `gorevolt_test_token` environment variable")
	}

	client := gorevolt.New(token)
	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}

	client.Close()
}
