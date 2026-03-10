package handlers

import (
	"strings"

	"github.com/orchestra-mcp/plugin-bridge-discord/internal"
)

// ToolsHandler lists available MCP tools and commands.
type ToolsHandler struct{}

// NewToolsHandler creates a new tools handler.
func NewToolsHandler() *ToolsHandler { return &ToolsHandler{} }

func (h *ToolsHandler) Name() string                      { return "tools" }
func (h *ToolsHandler) MatchesPrefix(content string) bool { return strings.ToLower(strings.TrimSpace(content)) == "tools" }
func (h *ToolsHandler) MatchesSlash(name string) bool     { return name == "tools" }

// HandleMessage handles a prefix command message.
func (h *ToolsHandler) HandleMessage(msg *internal.MessageCreate, api internal.HandlerAPI) {
	h.doTools(msg.ChannelID, api)
}

// HandleInteraction handles a slash command interaction.
func (h *ToolsHandler) HandleInteraction(ix *internal.InteractionCreate, api internal.HandlerAPI) {
	channelID := api.ChannelID()
	if ix.Message != nil {
		channelID = ix.Message.ChannelID
	}
	h.doTools(channelID, api)
	api.RespondInteraction(ix.ID, ix.Token, internal.InteractionResponseDeferred, "", nil)
}

func (h *ToolsHandler) doTools(channelID string, api internal.HandlerAPI) {
	// List features from the MCP
	result, err := api.CallTool("list_features", map[string]any{})
	if err != nil {
		// Fallback: just show available commands
		desc := "**Available Commands:**\n"
		desc += "`!chat <prompt>` - Chat with Claude\n"
		desc += "`!mcp <tool> [args]` - Execute MCP tool\n"
		desc += "`!status [project]` - Project status\n"
		desc += "`!stop [session]` - Stop session\n"
		desc += "`!ping` - Health check\n"
		desc += "`!tools` - This help\n"
		api.SendToChannel(channelID, "", internal.InfoEmbed("Available Commands", desc))
		return
	}

	api.SendToChannel(channelID, "", internal.InfoEmbed("MCP Tools", internal.Truncate(result, internal.SafeEmbedDesc)))
}

// SlashDef returns the slash command definition.
func (h *ToolsHandler) SlashDef() *internal.SlashCommandDef {
	return &internal.SlashCommandDef{
		Name:        "tools",
		Description: "List available commands and MCP tools",
	}
}
