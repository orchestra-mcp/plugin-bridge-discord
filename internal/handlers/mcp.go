package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/orchestra-mcp/plugin-bridge-discord/internal"
)

// McpHandler executes MCP tools from Discord.
type McpHandler struct{}

// NewMcpHandler creates a new MCP handler.
func NewMcpHandler() *McpHandler { return &McpHandler{} }

func (h *McpHandler) Name() string                      { return "mcp" }
func (h *McpHandler) MatchesPrefix(content string) bool { return strings.HasPrefix(strings.ToLower(content), "mcp ") }
func (h *McpHandler) MatchesSlash(name string) bool     { return name == "mcp" }

// HandleMessage handles a prefix command message.
func (h *McpHandler) HandleMessage(msg *internal.MessageCreate, api internal.HandlerAPI) {
	parts := strings.SplitN(strings.TrimPrefix(msg.Content, "mcp "), " ", 2)
	toolName := parts[0]
	var argsJSON string
	if len(parts) > 1 {
		argsJSON = parts[1]
	}
	h.doMcp(msg.ChannelID, toolName, argsJSON, api)
}

// HandleInteraction handles a slash command interaction.
func (h *McpHandler) HandleInteraction(ix *internal.InteractionCreate, api internal.HandlerAPI) {
	var toolName, argsJSON string
	for _, opt := range ix.Data.Options {
		switch opt.Name {
		case "tool":
			toolName = opt.OptionString()
		case "args":
			argsJSON = opt.OptionString()
		}
	}
	api.RespondInteraction(ix.ID, ix.Token, internal.InteractionResponseDeferred, "", nil)
	channelID := api.ChannelID()
	if ix.Message != nil {
		channelID = ix.Message.ChannelID
	}
	h.doMcp(channelID, toolName, argsJSON, api)
}

func (h *McpHandler) doMcp(channelID, toolName, argsJSON string, api internal.HandlerAPI) {
	if toolName == "" {
		api.SendToChannel(channelID, "", internal.InfoEmbed("Usage", "`!mcp <tool> [json-args]`"))
		return
	}

	args := make(map[string]any)
	if argsJSON != "" {
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			api.SendToChannel(channelID, "", internal.ErrorEmbed("Invalid JSON", err.Error()))
			return
		}
	}

	api.SendToChannel(channelID, "", internal.InfoEmbed("Running", fmt.Sprintf("`%s`", toolName)))

	result, err := api.CallTool(toolName, args)
	if err != nil {
		api.SendToChannel(channelID, "", internal.ErrorEmbed("Tool Error", err.Error()))
		return
	}

	if len(result) <= internal.SafeEmbedDesc {
		api.SendToChannel(channelID, "", internal.SuccessEmbed(fmt.Sprintf("%s", toolName), result))
	} else {
		api.SendToChannel(channelID, "", internal.SuccessEmbed(fmt.Sprintf("%s", toolName), internal.Truncate(result, internal.SafeEmbedDesc)))
	}
}

// SlashDef returns the slash command definition.
func (h *McpHandler) SlashDef() *internal.SlashCommandDef {
	return &internal.SlashCommandDef{
		Name:        "mcp",
		Description: "Execute an MCP tool",
		Options: []internal.SlashOptionDef{
			{Name: "tool", Description: "Tool name", Type: 3, Required: true},
			{Name: "args", Description: "JSON arguments", Type: 3, Required: false},
		},
	}
}
