package discordgo

import (
	"strconv"
	"testing"
	"time"
)

func TestChannel_SendMessage(t *testing.T) {
	if envChannel == "" {
		t.Skip("Skipping, DG_CHANNEL not set.")
	}

	if dg == nil {
		t.Skip("Skipping, dg not set.")
	}

	c, err := dg.State.Channel(envChannel)
	if err != nil {
		t.Fatalf("Channel %s wasn't cached", envChannel)
	}

	_, err = c.SendMessage("Testing Channel.SendMessage", nil, nil)
	if err != nil {
		t.Fatalf("Error while sending message: %s", err)
	}
}

func TestChannel_PermissionsFor(t *testing.T) {
	if envChannel == "" {
		t.Skip("Skipping, DG_CHANNEL not set.")
	}

	if envAdmin == "" {
		t.Skip("Skipping, DG_ADMIN not set.")
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

	m, err := g.GetMember(envAdmin)
	if err != nil {
		t.Fatalf("User %s is not in Guild", envAdmin)
	}

	p, err := c.PermissionsFor(m)
	if err != nil {
		t.Fatalf("Permissions calculation failed")
	}

	if !p.Has(PermissionAdministrator) {
		println(strconv.FormatInt(int64(p), 2))
		t.Fatalf("Envadmin does not have admin even though he should")
	}

	p.Set(PermissionManageNicknames, false)
	if p.Has(PermissionManageNicknames) {
		t.Fatalf("Permissions has manage nicknames even though it shouldn't")
	}

	p.Set(PermissionManageNicknames, true)
	if !p.Has(PermissionManageNicknames) {
		t.Fatalf("Permissions does not have manage nicknames even though it should")
	}
}

func TestChannel_MessagesBulkDelete(t *testing.T) {
	if envChannel == "" {
		t.Skip("Skipping, DG_CHANNEL not set.")
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

	msg1, err := c.SendMessage("Testing Channel.MessagesBulkDelete", nil, nil)
	if err != nil {
		t.Fatalf("Error while sending message: %s", err)
	}

	msg2, err := c.SendMessage("Testing Channel.MessagesBulkDelete", nil, nil)
	if err != nil {
		t.Fatalf("Error while sending message: %s", err)
	}

	done := make(chan bool, 1)

	f := dg.AddHandler(func(s *Session, event *MessageDeleteBulk) {
		if Contains(event.Messages, msg1.ID) && Contains(event.Messages, msg2.ID) {
			done <- true
		}
		t.Fatal("bulk delete did not delete the messages")
	})

	err = c.MessagesBulkDelete([]*Message{msg1, msg2})
	if err != nil {
		t.Fatal("error while doing bulk delete ", err)
	}

	select {
	case <-time.After(2000 * time.Millisecond):
		t.Fatal("no bulk delete happened")
	case <-done:
	}
	f()
}
