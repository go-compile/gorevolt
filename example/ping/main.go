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

	client.OnReady(func(c *gorevolt.Client, startup time.Duration) {
		fmt.Printf("[CONNECTED] [USER: %s]\n", c.User.Username)
	})

	client.OnMessage(func(c *gorevolt.Client, m *gorevolt.Message) {
		// ignore self
		if m.AuthorID == c.User.ID {
			return
		}

		fmt.Printf("[NEW MESSAGE] [USER: %s] [SERVER: %s] [CHANNEL: %s] %q\n", m.AuthorID, m.Server().Name, m.Channel.Name, m.Content)
		m.Reply("Hello there " + m.Author.Username)
	})

	client.OnMessageUpdate(func(c *gorevolt.Client, m *gorevolt.UpdatedMessage) {
		fmt.Println("Message updated:", m.Edited.Day())
		fmt.Println("Message:", m.Content)
	})

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	<-close
	client.Close()
}
