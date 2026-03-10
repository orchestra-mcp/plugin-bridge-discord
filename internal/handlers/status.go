package handlers

import (
	"strings"

	"github.com/orchestra-mcp/plugin-bridge-discord/internal"
)

// StatusHandler shows project workflow status.
type StatusHandler struct{}

// NewStatusHandler creates a new status handler.
func NewStatusHandler() *StatusHandler { return &StatusHandler{} }

func (h *StatusHandler) Name() string                      { return "status" }
func (h *StatusHandler) MatchesPrefix(content string) bool { return strings.HasPrefix(strings.ToLower(content), "status") }
func (h *StatusHandler) MatchesSlash(name string) bool     { return name == "status" }

// HandleMessage handles a prefix command message.
func (h *StatusHandler) HandleMessage(msg *internal.MessageCreate, api internal.HandlerAPI) {
	parts := strings.Fields(msg.Content)
	projectID := ""
	if len(parts) > 1 {
		projectID = parts[1]
	}
	h.doStatus(msg.ChannelID, projectID, api)
}

// HandleInteraction handles a slash command interaction.
func (h *StatusHandler) HandleInteraction(ix *internal.InteractionCreate, api internal.HandlerAPI) {
	var projectID string
	for _, opt := range ix.Data.Options {
		if opt.Name == "project" {
			projectID = opt.OptionString()
		}
	}
	api.RespondInteraction(ix.ID, ix.Token, internal.InteractionResponseDeferred, "", nil)
	channelID := api.ChannelID()
	if ix.Message != nil {
		channelID = ix.Message.ChannelID
	}
	h.doStatus(channelID, projectID, api)
}

func (h *StatusHandler) doStatus(channelID, projectID string, api internal.HandlerAPI) {
	args := map[string]any{}
	if projectID != "" {
		args["project_id"] = projectID
	}

	result, err := api.CallTool("get_project_status", args)
	if err != nil {
		// Try get_progress as fallback
		result, err = api.CallTool("get_progress", args)
		if err != nil {
			api.SendToChannel(channelID, "", internal.ErrorEmbed("Status Error", err.Error()))
			return
		}
	}

	api.SendToChannel(channelID, "", internal.InfoEmbed("Project Status", result))
}

// SlashDef returns the slash command definition.
func (h *StatusHandler) SlashDef() *internal.SlashCommandDef {
	return &internal.SlashCommandDef{
		Name:        "status",
		Description: "Show project workflow status",
		Options: []internal.SlashOptionDef{
			{Name: "project", Description: "Project slug", Type: 3, Required: false},
		},
	}
}
