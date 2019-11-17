package discordgo

import (
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

//////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////// VARS NEEDED FOR TESTING
var (
	dg *Session // Stores a global discordgo bot session

	envToken     = os.Getenv("DG_TOKEN")      // Token to use when authenticating the bot account
	envGuild     = os.Getenv("DG_GUILD")      // Guild ID to use for tests
	envChannel   = os.Getenv("DG_CHANNEL")    // Channel ID to use for tests
	envAdmin     = os.Getenv("DG_ADMIN")      // User ID of admin user to use for tests
	envRole      = os.Getenv("DG_ROLE")       // Role ID of role to use for tests
	envAdminRole = os.Getenv("DG_ADMIN_ROLE") // Role ID of the highest role in use
)

func init() {
	fmt.Println("Init is being called.")
	if envToken != "" {
		if d, err := New(envToken); err == nil {
			dg = d
		}
	}
}

//////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////// START OF TESTS

// TestNew tests the New() function without any arguments.  This should return
// a valid Session{} struct and no errors.
func TestNew(t *testing.T) {

	_, err := New()
	if err != nil {
		t.Errorf("New() returned error: %+v", err)
	}
}

// TestInvalidToken tests the New() function with an invalid token
func TestInvalidToken(t *testing.T) {
	_, err := New("asjkldhflkjasdh")
	if err != nil {
		t.Fatalf("New(InvalidToken) returned error: %+v", err)
	}
}

// TestNewToken tests the New() function with a Token.
func TestNewToken(t *testing.T) {

	if envToken == "" {
		t.Skip("Skipping New(token), DG_TOKEN not set")
	}

	d, err := New(envToken)
	if err != nil {
		t.Fatalf("New(envToken) returned error: %+v", err)
	}

	if d == nil {
		t.Fatal("New(envToken), d is nil, should be Session{}")
	}

	if d.Token == "" {
		t.Fatal("New(envToken), d.Token is empty, should be a valid Token.")
	}
}

func TestOpenClose(t *testing.T) {
	if envToken == "" {
		t.Skip("Skipping TestClose, DG_TOKEN not set")
	}

	d, err := New(envToken)
	if err != nil {
		t.Fatalf("TestClose, New(envToken) returned error: %+v", err)
	}

	if err = d.Open(); err != nil {
		t.Fatalf("TestClose, d.Open failed: %+v", err)
	}

	// We need a better way to know the session is ready for use,
	// this is totally gross.
	start := time.Now()
	for {
		d.RLock()
		if d.DataReady {
			d.RUnlock()
			break
		}
		d.RUnlock()

		if time.Since(start) > 10*time.Second {
			t.Fatal("DataReady never became true.yy")
		}
		runtime.Gosched()
	}

	// TODO find a better way
	// Add a small sleep here to make sure heartbeat and other events
	// have enough time to get fired.  Need a way to actually check
	// those events.
	time.Sleep(2 * time.Second)

	// UpdateStatus - maybe we move this into wsapi_test.go but the websocket
	// created here is needed.  This helps tests that the websocket was setup
	// and it is working.
	if err = d.UpdateStatus(0, time.Now().String()); err != nil {
		t.Errorf("UpdateStatus error: %+v", err)
	}

	if err = d.Close(); err != nil {
		t.Fatalf("TestClose, d.Close failed: %+v", err)
	}

	if dg.State == nil {
		t.Fatal("state is nil")
	}

	done := make(chan bool, 1)
	f := dg.AddHandler(func(s *Session, _ *Ready) {
		done <- true
	})

	err = dg.Open()
	if err != nil {
		t.Fatal("Opening of actual connection with discord failed")
	}

	select {
	case <-time.After(2000 * time.Millisecond):
		t.Fatal("didn't receive ready")
	case <-done:
	}
	f()
}

func TestAddHandler(t *testing.T) {

	testHandlerCalled := int32(0)
	testHandler := func(s *Session, m *MessageCreate) {
		atomic.AddInt32(&testHandlerCalled, 1)
	}

	interfaceHandlerCalled := int32(0)
	interfaceHandler := func(s *Session, i interface{}) {
		atomic.AddInt32(&interfaceHandlerCalled, 1)
	}

	bogusHandlerCalled := int32(0)
	bogusHandler := func(s *Session, se *Session) {
		atomic.AddInt32(&bogusHandlerCalled, 1)
	}

	d := Session{}
	d.AddHandler(testHandler)
	d.AddHandler(testHandler)

	d.AddHandler(interfaceHandler)
	d.AddHandler(bogusHandler)

	d.handleEvent(messageCreateEventType, &MessageCreate{})
	d.handleEvent(messageDeleteEventType, &MessageDelete{})

	<-time.After(500 * time.Millisecond)

	// testHandler will be called twice because it was added twice.
	if atomic.LoadInt32(&testHandlerCalled) != 2 {
		t.Fatalf("testHandler was not called twice.")
	}

	// interfaceHandler will be called twice, once for each event.
	if atomic.LoadInt32(&interfaceHandlerCalled) != 2 {
		t.Fatalf("interfaceHandler was not called twice.")
	}

	if atomic.LoadInt32(&bogusHandlerCalled) != 0 {
		t.Fatalf("bogusHandler was called.")
	}
}

func TestRemoveHandler(t *testing.T) {

	testHandlerCalled := int32(0)
	testHandler := func(s *Session, m *MessageCreate) {
		atomic.AddInt32(&testHandlerCalled, 1)
	}

	d := Session{}
	r := d.AddHandler(testHandler)

	d.handleEvent(messageCreateEventType, &MessageCreate{})

	r()

	d.handleEvent(messageCreateEventType, &MessageCreate{})

	<-time.After(500 * time.Millisecond)

	// testHandler will be called once, as it was removed in between calls.
	if atomic.LoadInt32(&testHandlerCalled) != 1 {
		t.Fatalf("testHandler was not called once.")
	}
}

func TestHandlerSessionInserter(t *testing.T) {
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

	done := make(chan bool, 1)

	testChannelHandler := func(s *Session, c *ChannelCreate) {
		_, _ = c.SendMessage("OwO A new channel was made", nil, nil)
	}

	testMessageHandler := func(s *Session, m *MessageCreate) {
		_, _ = m.Edit(m.NewMessageEdit().SetContent("OwO message received and edited"))
	}

	testMessageUpdateHandler := func(s *Session, m *MessageUpdate) {
		_ = m.AddReaction(&Emoji{Name: "❤"})
		done <- true
	}

	r := dg.AddHandler(testChannelHandler)
	m := dg.AddHandler(testMessageHandler)
	u := dg.AddHandler(testMessageUpdateHandler)

	_, err = g.CreateChannel("TestChannel", ChannelTypeGuildText)
	if err != nil {
		r()
		return
	}

	select {
	case <-time.After(2000 * time.Millisecond):
		t.Fatal("the handlers weren't called")
	case <-done:
	}
	r()
	m()
	u()
}
