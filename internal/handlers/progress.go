package handlers

import (
	"strings"

	"github.com/orchestra-mcp/plugin-bridge-discord/internal"
)

// ProgressHandler watches session progress.
type ProgressHandler struct{}

// NewProgressHandler creates a new progress handler.
func NewProgressHandler() *ProgressHandler { return &ProgressHandler{} }

func (h *ProgressHandler) Name() string                      { return "watch" }
func (h *ProgressHandler) MatchesPrefix(content string) bool { return strings.HasPrefix(strings.ToLower(content), "watch") }
func (h *ProgressHandler) MatchesSlash(name string) bool     { return name == "watch" }

// HandleMessage handles a prefix command message.
func (h *ProgressHandler) HandleMessage(msg *internal.MessageCreate, api internal.HandlerAPI) {
	parts := strings.Fields(msg.Content)
	sessionID := ""
	if len(parts) > 1 {
		sessionID = parts[1]
	}
	h.doWatch(msg.ChannelID, sessionID, api)
}

// HandleInteraction handles a slash command interaction.
func (h *ProgressHandler) HandleInteraction(ix *internal.InteractionCreate, api internal.HandlerAPI) {
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
	h.doWatch(channelID, sessionID, api)
}

func (h *ProgressHandler) doWatch(channelID, sessionID string, api internal.HandlerAPI) {
	if sessionID == "" {
		result, err := api.CallTool("list_active", map[string]any{})
		if err != nil {
			api.SendToChannel(channelID, "", internal.ErrorEmbed("Error", err.Error()))
			return
		}
		api.SendToChannel(channelID, "", internal.InfoEmbed("Active Sessions", result+"\n\nUse `!watch <session_id>` to watch"))
		return
	}

	result, err := api.CallTool("session_status", map[string]any{"session_id": sessionID})
	if err != nil {
		api.SendToChannel(channelID, "", internal.ErrorEmbed("Watch Error", err.Error()))
		return
	}
	api.SendToChannel(channelID, "", internal.InfoEmbed("Session Progress", result))
}

// SlashDef returns the slash command definition.
func (h *ProgressHandler) SlashDef() *internal.SlashCommandDef {
	return &internal.SlashCommandDef{
		Name:        "watch",
		Description: "Watch session progress",
		Options: []internal.SlashOptionDef{
			{Name: "session", Description: "Session ID to watch", Type: 3, Required: false},
		},
	}
}
