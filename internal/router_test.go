package internal

import (
	"sync"
	"testing"
)

// mockHandler is a test handler for router tests.
type mockHandler struct {
	name      string
	prefix    string
	slash     string
	called    bool
	mu        sync.Mutex
}

func (h *mockHandler) Name() string { return h.name }
func (h *mockHandler) MatchesPrefix(content string) bool {
	return len(content) >= len(h.prefix) && content[:len(h.prefix)] == h.prefix
}
func (h *mockHandler) MatchesSlash(name string) bool { return name == h.slash }
func (h *mockHandler) HandleMessage(msg *MessageCreate, api HandlerAPI) {
	h.mu.Lock()
	h.called = true
	h.mu.Unlock()
}
func (h *mockHandler) HandleInteraction(ix *InteractionCreate, api HandlerAPI) {
	h.mu.Lock()
	h.called = true
	h.mu.Unlock()
}
func (h *mockHandler) SlashDef() *SlashCommandDef {
	if h.slash == "" {
		return nil
	}
	return &SlashCommandDef{Name: h.slash, Description: "test"}
}

// mockAPI implements HandlerAPI for testing.
type mockAPI struct {
	cfg *Config
}

func (m *mockAPI) SendToChannel(channelID, content string, embed *Embed) error { return nil }
func (m *mockAPI) SendComponents(channelID, content string, embed *Embed, components []Component) (string, error) {
	return "msg-1", nil
}
func (m *mockAPI) EditMessage(channelID, messageID, content string, embed *Embed) error { return nil }
func (m *mockAPI) RespondInteraction(id, token string, respType int, content string, embed *Embed) error {
	return nil
}
func (m *mockAPI) IsRunning() bool                              { return true }
func (m *mockAPI) ChannelID() string                            { return "ch-1" }
func (m *mockAPI) Config() *Config                              { return m.cfg }
func (m *mockAPI) CallTool(name string, args map[string]any) (string, error) { return "ok", nil }

func TestRouter_SlashDefs(t *testing.T) {
	r := NewRouter("!")
	r.Register(&mockHandler{name: "ping", slash: "ping"})
	r.Register(&mockHandler{name: "no-slash"})

	defs := r.SlashDefs()
	if len(defs) != 1 {
		t.Errorf("SlashDefs count = %d, want 1", len(defs))
	}
	if defs[0].Name != "ping" {
		t.Errorf("SlashDef name = %q, want ping", defs[0].Name)
	}
}

func TestRouter_DefaultPrefix(t *testing.T) {
	r := NewRouter("")
	if r.prefix != "!" {
		t.Errorf("default prefix = %q, want !", r.prefix)
	}
}

func TestRouter_RoutesBot(t *testing.T) {
	r := NewRouter("!")
	api := &mockAPI{cfg: DefaultConfig()}

	// Bot messages should be ignored
	r.RouteMessage(MessageCreate{
		Content: "!ping",
		Author:  Author{ID: "bot", Bot: true},
	}, api)
	// No crash = pass
}

func TestRouter_RoutesNonPrefix(t *testing.T) {
	r := NewRouter("!")
	h := &mockHandler{name: "ping", prefix: "ping"}
	r.Register(h)
	api := &mockAPI{cfg: DefaultConfig()}

	// Non-prefix messages should be ignored
	r.RouteMessage(MessageCreate{
		Content: "hello",
		Author:  Author{ID: "user1"},
	}, api)
	// Handler should not be called (async, but immediately returns)
}

func TestRouter_AllowedUsers(t *testing.T) {
	r := NewRouter("!")
	h := &mockHandler{name: "ping", prefix: "ping"}
	r.Register(h)
	api := &mockAPI{cfg: &Config{AllowedUsers: []string{"allowed-id"}}}

	// Disallowed user
	r.RouteMessage(MessageCreate{
		Content: "!ping",
		Author:  Author{ID: "not-allowed"},
	}, api)
	// Should not route (allowed users check)
}
