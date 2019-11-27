package discordgo

import (
	"log"
	"testing"
	"time"
)

func getGuild(t *testing.T) (g *Guild) {
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

	if g.Unavailable {
		t.Fatalf("Guild %s is still unavailable", envGuild)
	}
	return
}

func TestGuild_GetChannel(t *testing.T) {
	g := getGuild(t)

	_, err := g.GetChannel(envChannel)
	if err != nil {
		t.Fatalf("Channel not found in guild")
	}
}

func TestGuild_GetRole(t *testing.T) {
	g := getGuild(t)

	_, err := g.GetRole(envRole)
	if err != nil {
		t.Fatalf("Role not found in guild")
	}
}

func TestGuild_CreateDeleteRole(t *testing.T) {
	g := getGuild(t)

	r, err := g.CreateRole(&RoleSettings{})
	if err != nil {
		t.Fatalf("Role failed to create in Guild; %s", err)
	}

	editData := &RoleSettings{
		Name:        "OwO a testing role",
		Hoist:       false,
		Color:       ColorGreen,
		Permissions: r.Permissions,
		Mentionable: true,
	}

	r, err = r.EditComplex(editData)
	if err != nil {
		t.Fatalf("Failed at editing role; %s", err)
	}

	err = r.Delete()
	if err != nil {
		t.Fatalf("Failed at deleteing role; %s", err)
	}
}

func TestRole_Move(t *testing.T) {
	g := getGuild(t)

	r, err := g.CreateRole(&RoleSettings{})
	if err != nil {
		t.Fatalf("Role failed to create in Guild; %s", err)
	}

	c := ColorGreen
	err = c.SetHex("#ffff00")
	if err != nil {
		t.Fatalf("failed at parsing hex code; %s", err)
	}

	editData := &RoleSettings{
		Name:        "OwO a moving role",
		Hoist:       false,
		Color:       c,
		Permissions: r.Permissions,
		Mentionable: false,
	}

	r, err = r.EditComplex(editData)
	if err != nil {
		t.Fatalf("Failed at editing role; %s", err)
	}

	done := make(chan bool, 1)
	println(r.ID)
	testRoleUpdateHandler := func(se *Session, role *GuildRoleUpdate) {
		log.Println(role.Role.ID, role.BeforeUpdate.ID, role.Role.Position, role.BeforeUpdate.Position)

		if r.ID == role.Role.ID && role.Role.Position != role.BeforeUpdate.Position && role.Role.Position == 6 {
			done <- true
		}
	}
	rf := dg.AddHandler(testRoleUpdateHandler)

	err = r.Move(6)
	if err != nil {
		rf()
		t.Fatalf("failed at moving role; %s", err)
	}

	select {
	case <-time.After(2000 * time.Millisecond):
		rf()
		t.Fatal("the role handler wasn't called or the position wasn't changed")
	case <-done:
	}

	rf()

	err = r.Delete()
	if err != nil {
		t.Fatalf("Failed at deleteing role; %s", err)
	}
}
