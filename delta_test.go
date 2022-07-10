package gorevolt_test

import (
	"os"
	"testing"

	"github.com/go-compile/gorevolt"
)

var (
	testChannel = ""
	testUserID  = ""
)

func TestGetChannel(t *testing.T) {
	client := newTestingClient(t)

	channel, err := gorevolt.GetChannel(client, testChannel)
	if err != nil {
		t.Fatal(err)
	}

	if channel.ChannelType != "TextChannel" {
		t.Error("channel type incorrect")
	}

	if channel.ID != testChannel {
		t.Error("channel ID incorrect")
	}
}

func TestGetUser(t *testing.T) {
	client := newTestingClient(t)

	user, err := gorevolt.GetUser(client, testUserID)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != testUserID {
		t.Error("user ID incorrect")
	}

	if user.Bot.OwnerID == "" {
		t.Error("bot owner ID not set")
	}

	if user.Username != "GoRevolt" {
		t.Errorf("bot name %q incorrect should be GoRevolt\n", user.Username)
	}
}

func newTestingClient(t *testing.T) *gorevolt.Client {
	token, exists := os.LookupEnv("gorevolt_test_token")
	if !exists {
		t.Skip("no token provided in `gorevolt_test_token` environment variable")
	}

	channelID, exists := os.LookupEnv("gorevolt_test_channel")
	if !exists {
		t.Skip("no token provided in `gorevolt_test_channel` environment variable")
	}

	userID, exists := os.LookupEnv("gorevolt_test_user")
	if !exists {
		t.Skip("no token provided in `gorevolt_test_user` environment variable")
	}

	testChannel = channelID
	testUserID = userID

	return gorevolt.New(token)
}
