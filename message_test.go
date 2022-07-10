package gorevolt_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-compile/gorevolt"
)

func TestMessageHandler(t *testing.T) {
	token, exists := os.LookupEnv("gorevolt_test_token")
	if !exists {
		t.Skip("no token provided in `gorevolt_test_token` environment variable")
	}

	interactiveTests, _ := os.LookupEnv("gorevolt_test_interactive")
	if interactiveTests != "true" {
		t.Skip("`gorevolt_test_interactive` environment variable needs to be set to true")

	}

	client := gorevolt.New(token)
	close := make(chan struct{})

	client.OnMessage(func(c *gorevolt.Client, m *gorevolt.Message) {
		fmt.Printf("[NEW MESSAGE] [USER: %s] [SERVER: %s] [CHANNEL: %s] %q\n", m.AuthorID, m.Server().Name, m.Channel.Name, m.Content)

		fmt.Println(m.Author.Username)

		// close and finish test
		close <- struct{}{}
	})

	if err := client.Connect(); err != nil {
		t.Fatal(err)
	}

	fmt.Println("Send a message in a channel which the Revolt has access to")

	<-close
	client.Close()
}
