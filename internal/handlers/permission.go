package handlers

import (
	"strings"

	"github.com/orchestra-mcp/plugin-bridge-discord/internal"
)

// PermissionHandler handles tool permission approve/deny via Discord buttons.
type PermissionHandler struct{}

// NewPermissionHandler creates a new permission handler.
func NewPermissionHandler() *PermissionHandler { return &PermissionHandler{} }

func (h *PermissionHandler) Name() string                                                        { return "permission" }
func (h *PermissionHandler) MatchesPrefix(_ string) bool                                         { return false }
func (h *PermissionHandler) MatchesSlash(_ string) bool                                          { return false }
func (h *PermissionHandler) HandleMessage(_ *internal.MessageCreate, _ internal.HandlerAPI)       {}
func (h *PermissionHandler) HandleInteraction(_ *internal.InteractionCreate, _ internal.HandlerAPI) {}
func (h *PermissionHandler) SlashDef() *internal.SlashCommandDef                                 { return nil }

// MatchesCustomID matches permission button interactions.
func (h *PermissionHandler) MatchesCustomID(customID string) bool {
	return strings.HasPrefix(customID, "perm_approve_") || strings.HasPrefix(customID, "perm_deny_")
}

// HandleComponentInteraction handles permission button clicks.
func (h *PermissionHandler) HandleComponentInteraction(ix *internal.InteractionCreate, api internal.HandlerAPI) {
	customID := ix.Data.CustomID
	var decision, reqID string

	if strings.HasPrefix(customID, "perm_approve_") {
		decision = "approve"
		reqID = strings.TrimPrefix(customID, "perm_approve_")
	} else if strings.HasPrefix(customID, "perm_deny_") {
		decision = "deny"
		reqID = strings.TrimPrefix(customID, "perm_deny_")
	}

	_, err := api.CallTool("respond_permission", map[string]any{
		"id":       reqID,
		"decision": decision,
	})
	if err != nil {
		api.RespondInteraction(ix.ID, ix.Token, internal.InteractionResponseUpdate, "",
			internal.ErrorEmbed("Permission Error", err.Error()))
		return
	}

	status := "Approved"
	if decision == "deny" {
		status = "Denied"
	}
	api.RespondInteraction(ix.ID, ix.Token, internal.InteractionResponseUpdate, "",
		internal.PermissionResultEmbed(status, reqID))
}
