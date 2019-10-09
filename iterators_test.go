package discordgo

import (
	"testing"
)

func TestHistoryIterator_Iter(t *testing.T) {
	if envChannel == "" {
		t.Skip("Skipping, DG_CHANNEL not set.")
	}

	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	g, err := dg.State.Guild(envGuild)
	if err != nil {
		t.Fatalf("Guild not found, id: %s; %s", envGuild, err)
	}

	c, err := g.GetChannel(envChannel)
	if err != nil || c == nil {
		t.Fatalf("Channel %s wasn't cached", envChannel)
	}

	iterator := c.GetHistoryIterator().SetLimit(10)
	messages := make([]*Message, 0, 10)
	for m := range iterator.Iter() {
		messages = append(messages, m)
	}

	if len(messages) != 10 {
		t.Fatal("iterator does not return 10 messages")
	}
}

func TestHistoryIterator_Reverse(t *testing.T) {
	if envChannel == "" {
		t.Skip("Skipping, DG_CHANNEL not set.")
	}

	if envGuild == "" {
		t.Skip("Skipping, DG_GUILD not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	g, err := dg.State.Guild(envGuild)
	if err != nil {
		t.Fatalf("Guild not found, id: %s; %s", envGuild, err)
	}

	c, err := g.GetChannel(envChannel)
	if err != nil || c == nil {
		t.Fatalf("Channel %s wasn't cached", envChannel)
	}

	iterator := c.GetHistoryIterator().SetLimit(10)
	messages := make([]*Message, 0, 10)
	for m := range iterator.Iter() {
		messages = append(messages, m)
	}

	if len(messages) != 10 {
		t.Fatal("iterator does not return 10 messages")
	}

	reverseMessages := make([]*Message, 0, 10)
	for m := range iterator.Reset().SetLimit(10).Reverse().Iter() {
		reverseMessages = append(reverseMessages, m)
	}

	if len(reverseMessages) != 10 {
		t.Fatal("reversed iterator does not return 10 messages")
	}

	if messages[0].ID != reverseMessages[len(reverseMessages)-1].ID {
		t.Fatal("message order not reversed")
	}
}
