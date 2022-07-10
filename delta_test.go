package gorevolt_test

import (
	"os"
	"testing"

	"github.com/go-compile/gorevolt"
)

var testChannel = ""

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

func newTestingClient(t *testing.T) *gorevolt.Client {
	token, exists := os.LookupEnv("gorevolt_test_token")
	if !exists {
		t.Skip("no token provided in `gorevolt_test_token` environment variable")
	}

	channelID, exists := os.LookupEnv("gorevolt_test_channel")
	if !exists {
		t.Skip("no token provided in `gorevolt_test_token` environment variable")
	}

	testChannel = channelID

	return gorevolt.New(token)
}
