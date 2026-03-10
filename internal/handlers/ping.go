package handlers

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/orchestra-mcp/plugin-bridge-discord/internal"
)

var startTime = time.Now()

// PingHandler provides a health check command.
type PingHandler struct{}

// NewPingHandler creates a new ping handler.
func NewPingHandler() *PingHandler { return &PingHandler{} }

func (h *PingHandler) Name() string                      { return "ping" }
func (h *PingHandler) MatchesPrefix(content string) bool { return strings.ToLower(strings.TrimSpace(content)) == "ping" }
func (h *PingHandler) MatchesSlash(name string) bool     { return name == "ping" }

// HandleMessage handles a prefix command message.
func (h *PingHandler) HandleMessage(msg *internal.MessageCreate, api internal.HandlerAPI) {
	h.doPing(msg.ChannelID, api)
}

// HandleInteraction handles a slash command interaction.
func (h *PingHandler) HandleInteraction(ix *internal.InteractionCreate, api internal.HandlerAPI) {
	channelID := api.ChannelID()
	if ix.Message != nil {
		channelID = ix.Message.ChannelID
	}
	h.doPing(channelID, api)
	api.RespondInteraction(ix.ID, ix.Token, internal.InteractionResponseMessage, "Pong!", nil)
}

func (h *PingHandler) doPing(channelID string, api internal.HandlerAPI) {
	uptime := time.Since(startTime).Round(time.Second)
	desc := fmt.Sprintf("**Uptime:** %s\n**Go:** %s\n**Platform:** %s/%s", uptime, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	api.SendToChannel(channelID, "", internal.SuccessEmbed("Pong!", desc))
}

// SlashDef returns the slash command definition.
func (h *PingHandler) SlashDef() *internal.SlashCommandDef {
	return &internal.SlashCommandDef{
		Name:        "ping",
		Description: "Health check",
	}
}
