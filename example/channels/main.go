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
		m.Reply("Hello there " + m.Channel.Name)
	})

	client.OnMessageUpdate(func(c *gorevolt.Client, m *gorevolt.UpdatedMessage) {
		fmt.Println("Message updated:", m.Edited.Day())
		fmt.Println("Message:", m.Content)
	})

	client.OnChannelCreate(func(c *gorevolt.Client, channel *gorevolt.Channel) {
		fmt.Printf("\nNew channel created [%s] in [%s]\n\n", channel.Name, channel.Server(c).Name)

		channel.SendMessagef(c, "Welcome to the new channel %s within %s.", channel.Name, channel.Server(c).Name)
	})

	client.OnChannelUpdate(func(c *gorevolt.Client, old, new *gorevolt.Channel) {
		new.SendMessagef(c, "Channel updated. Name: %s, Description: %s. \nOld name: %s", new.Name, new.Description, old.Name)
	})

	client.OnChannelDelete(func(c *gorevolt.Client, channel *gorevolt.Channel) {
		fmt.Printf("Channel %q deleted.\n", channel.Name)
	})

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	<-close
	client.Close()
}
