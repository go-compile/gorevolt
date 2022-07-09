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

	client.OnServerUpdate(func(c *gorevolt.Client, old, new *gorevolt.Server) {

		m := fmt.Sprintf("| Name        | %.20s | %.20s |\n", old.Name, new.Name)
		m += fmt.Sprintf("| Description | %.20s | %.20s |\n", old.Description, new.Description)
		m += fmt.Sprintf("| Channels    | %20.d | %20.d |\n", len(old.ChannelIDs), len(new.ChannelIDs))
		m += fmt.Sprintf("| Categories  | %20.d | %20.d |\n", len(old.Categories), len(new.Categories))
		fmt.Printf("Server %q updated.\n\n%s", old.Name, m)

		fmt.Println(old.Categories)
		fmt.Println(new.Categories)

	})

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	<-close
	client.Close()
}
