package handlers

import (
	"strings"

	"github.com/orchestra-mcp/plugin-bridge-discord/internal"
)

// StopHandler stops Claude sessions.
type StopHandler struct{}

// NewStopHandler creates a new stop handler.
func NewStopHandler() *StopHandler { return &StopHandler{} }

func (h *StopHandler) Name() string                      { return "stop" }
func (h *StopHandler) MatchesPrefix(content string) bool { return strings.HasPrefix(strings.ToLower(content), "stop") }
func (h *StopHandler) MatchesSlash(name string) bool     { return name == "stop" }

// HandleMessage handles a prefix command message.
func (h *StopHandler) HandleMessage(msg *internal.MessageCreate, api internal.HandlerAPI) {
	parts := strings.Fields(msg.Content)
	sessionID := ""
	if len(parts) > 1 {
		sessionID = parts[1]
	}
	h.doStop(msg.ChannelID, sessionID, api)
}

// HandleInteraction handles a slash command interaction.
func (h *StopHandler) HandleInteraction(ix *internal.InteractionCreate, api internal.HandlerAPI) {
	var sessionID string
	for _, opt := range ix.Data.Options {
		if opt.Name == "session" {
			sessionID = opt.OptionString()
		}
	}
	api.RespondInteraction(ix.ID, ix.Token, internal.InteractionResponseDeferred, "", nil)
	channelID := api.ChannelID()
	if ix.Message != nil {
		channelID = ix.Message.ChannelID
	}
	h.doStop(channelID, sessionID, api)
}

func (h *StopHandler) doStop(channelID, sessionID string, api internal.HandlerAPI) {
	if sessionID == "" {
		// List active sessions first
		result, err := api.CallTool("list_active", map[string]any{})
		if err != nil {
			api.SendToChannel(channelID, "", internal.ErrorEmbed("Error", err.Error()))
			return
		}
		api.SendToChannel(channelID, "", internal.InfoEmbed("Active Sessions", result+"\n\nUse `!stop <session_id>` to stop a session"))
		return
	}

	result, err := api.CallTool("kill_session", map[string]any{"session_id": sessionID})
	if err != nil {
		api.SendToChannel(channelID, "", internal.ErrorEmbed("Stop Error", err.Error()))
		return
	}
	api.SendToChannel(channelID, "", internal.SuccessEmbed("Session Stopped", result))
}

// SlashDef returns the slash command definition.
func (h *StopHandler) SlashDef() *internal.SlashCommandDef {
	return &internal.SlashCommandDef{
		Name:        "stop",
		Description: "Stop a Claude session",
		Options: []internal.SlashOptionDef{
			{Name: "session", Description: "Session ID to stop", Type: 3, Required: false},
		},
	}
}
