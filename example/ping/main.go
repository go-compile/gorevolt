package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-compile/gorevolt"
)

func main() {
	token, exists := os.LookupEnv("revolt_token")
	if !exists {
		log.Fatal("set your token as the ENV revolt_token")
	}

	client := gorevolt.New(token)
	close := make(chan struct{})

	client.Register(func(c *gorevolt.Client, startup time.Duration) {
		fmt.Printf("[CONNECTED] [USER: %s]\n", c.User.Username)
	})

	client.Register(func(c *gorevolt.Client, m *gorevolt.Message) {
		// ignore self
		if m.AuthorID == c.User.ID {
			return
		}

		fmt.Printf("[NEW MESSAGE] [USER: %s] [SERVER: %s] [CHANNEL: %s] %q\n", m.AuthorID, m.Server().Name, m.Channel.Name, m.Content)

		author, err := m.Author()
		if err != nil {
			log.Println(err)
		}

		m.Reply("Hello there " + author.Username)
	})

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	<-close
	client.Close()
}
