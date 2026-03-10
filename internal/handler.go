package internal

// HandlerAPI provides capabilities to handlers for sending messages
// and accessing shared services.
type HandlerAPI interface {
	SendToChannel(channelID, content string, embed *Embed) error
	SendComponents(channelID, content string, embed *Embed, components []Component) (string, error)
	EditMessage(channelID, messageID, content string, embed *Embed) error
	RespondInteraction(id, token string, respType int, content string, embed *Embed) error
	IsRunning() bool
	ChannelID() string
	Config() *Config
	// CallTool invokes an MCP tool by name via cross-plugin call.
	CallTool(name string, args map[string]any) (string, error)
}

// Handler processes Discord commands (both prefix and slash).
type Handler interface {
	Name() string
	MatchesPrefix(content string) bool
	MatchesSlash(name string) bool
	HandleMessage(msg *MessageCreate, api HandlerAPI)
	HandleInteraction(ix *InteractionCreate, api HandlerAPI)
	SlashDef() *SlashCommandDef
}

// InteractionHandler is an optional interface for handlers that process
// button/component interactions (INTERACTION_CREATE type 3).
type InteractionHandler interface {
	MatchesCustomID(customID string) bool
	HandleComponentInteraction(ix *InteractionCreate, api HandlerAPI)
}
