package tools

import (
	"github.com/orchestra-mcp/plugin-bridge-discord/internal"
)

// DiscordBridge provides access to the Discord bot for MCP tools.
type DiscordBridge struct {
	Plugin *internal.BridgePlugin
}
