package handlers

import (
	"fmt"
	"strings"
	"sync"

	"github.com/orchestra-mcp/plugin-bridge-discord/internal"
)

// ChatHandler handles Claude chat via Discord with sticky channel-session mapping.
type ChatHandler struct {
	mu             sync.Mutex
	channelSession map[string]string // channelID -> session ID
}

// NewChatHandler creates a new chat handler.
func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		channelSession: make(map[string]string),
	}
}

func (h *ChatHandler) Name() string                      { return "chat" }
func (h *ChatHandler) MatchesPrefix(content string) bool { return strings.HasPrefix(strings.ToLower(content), "chat ") }
func (h *ChatHandler) MatchesSlash(name string) bool     { return name == "chat" }

// HandleMessage handles a prefix command message.
func (h *ChatHandler) HandleMessage(msg *internal.MessageCreate, api internal.HandlerAPI) {
	prompt := strings.TrimPrefix(msg.Content, "chat ")
	prompt = strings.TrimPrefix(prompt, "Chat ")
	if prompt == "" {
		api.SendToChannel(msg.ChannelID, "", internal.InfoEmbed("Usage", "`!chat <prompt>`"))
		return
	}
	h.doChat(msg.ChannelID, prompt, api)
}

// HandleInteraction handles a slash command interaction.
func (h *ChatHandler) HandleInteraction(ix *internal.InteractionCreate, api internal.HandlerAPI) {
	var prompt string
	for _, opt := range ix.Data.Options {
		if opt.Name == "prompt" {
			prompt = opt.OptionString()
		}
	}
	if prompt == "" {
		api.RespondInteraction(ix.ID, ix.Token, internal.InteractionResponseMessage, "Please provide a prompt", nil)
		return
	}
	api.RespondInteraction(ix.ID, ix.Token, internal.InteractionResponseDeferred, "", nil)

	channelID := ""
	if ix.Message != nil {
		channelID = ix.Message.ChannelID
	}
	if channelID == "" {
		channelID = api.ChannelID()
	}
	h.doChat(channelID, prompt, api)
}

func (h *ChatHandler) doChat(channelID, prompt string, api internal.HandlerAPI) {
	api.SendToChannel(channelID, "", internal.InfoEmbed("Processing", fmt.Sprintf("```\n%s\n```", internal.Truncate(prompt, 200))))

	// Call ai_prompt tool via cross-plugin bridge
	result, err := api.CallTool("ai_prompt", map[string]any{
		"prompt": prompt,
		"wait":   true,
	})
	if err != nil {
		api.SendToChannel(channelID, "", internal.ErrorEmbed("Error", err.Error()))
		return
	}

	// Split long responses for Discord's 4096 char embed limit
	if len(result) <= internal.SafeEmbedDesc {
		api.SendToChannel(channelID, "", internal.SuccessEmbed("Response", result))
		return
	}
	chunks := splitMessage(result, internal.SafeEmbedDesc)
	for i, chunk := range chunks {
		title := fmt.Sprintf("Response (%d/%d)", i+1, len(chunks))
		api.SendToChannel(channelID, "", internal.SuccessEmbed(title, chunk))
	}
}

// SlashDef returns the slash command definition.
func (h *ChatHandler) SlashDef() *internal.SlashCommandDef {
	return &internal.SlashCommandDef{
		Name:        "chat",
		Description: "Chat with Claude via Orchestra",
		Options: []internal.SlashOptionDef{
			{Name: "prompt", Description: "Your prompt", Type: 3, Required: true},
		},
	}
}

func splitMessage(s string, maxLen int) []string {
	var chunks []string
	for len(s) > maxLen {
		chunks = append(chunks, s[:maxLen])
		s = s[maxLen:]
	}
	if len(s) > 0 {
		chunks = append(chunks, s)
	}
	return chunks
}
