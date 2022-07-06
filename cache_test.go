package gorevolt_test

import (
	"strconv"
	"testing"

	"github.com/go-compile/gorevolt"
)

func TestArrayCacheUser(t *testing.T) {
	c := gorevolt.NewArrayCache(101, 10)

	for i := 1; i <= 100; i++ {
		c.PutUser(&gorevolt.User{
			ID:       strconv.Itoa(i),
			Username: "User-" + strconv.Itoa(i),
		})
	}

	user := &gorevolt.User{
		ID:       "101",
		Username: "BotUser",
	}

	c.PutUser(user)

	u := c.GetUser(user.ID)
	if u == nil {
		t.Fatal("failed to fetch user")
	}

	if u.ID != user.ID || u.Username != user.Username {
		t.Fatal("wrong user fetched")
	}
}

func TestArrayCacheEditUser(t *testing.T) {
	c := gorevolt.NewArrayCache(101, 10)

	user := &gorevolt.User{
		ID:       "101",
		Username: "BotUser",
	}

	c.PutUser(user)

	for i := 1; i <= 100; i++ {
		c.PutUser(&gorevolt.User{
			ID:       strconv.Itoa(i),
			Username: "User-" + strconv.Itoa(i),
		})
	}

	user.Username = "NewUsername"

	c.PutUser(user)

	u := c.GetUser(user.ID)
	if u == nil {
		t.Fatal("failed to fetch user")
	}
}
