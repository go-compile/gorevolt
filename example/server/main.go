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

	client.OnServerCreate(func(c *gorevolt.Client, server *gorevolt.Server) {
		fmt.Printf("Server created/joined: %s\n", server.Name)
	})

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	<-close
	client.Close()
}
