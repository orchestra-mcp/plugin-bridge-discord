package tools

import (
	"github.com/orchestra-mcp/plugin-bridge-discord/internal"
	"github.com/orchestra-mcp/sdk-go/plugin"
)

// RegisterAll registers all Discord bridge tools with the plugin builder.
func RegisterAll(builder *plugin.PluginBuilder, bp *internal.BridgePlugin) {
	bridge := &DiscordBridge{Plugin: bp}

	builder.RegisterTool("start_discord_bot",
		"Start the Discord bot (connects to gateway, registers commands)",
		StartBotSchema(), StartBot(bridge))

	builder.RegisterTool("stop_discord_bot",
		"Stop the Discord bot (disconnects from gateway)",
		StopBotSchema(), StopBot(bridge))

	builder.RegisterTool("discord_bot_status",
		"Get Discord bot status (running, config, allowed users)",
		BotStatusSchema(), BotStatus(bridge))

	builder.RegisterTool("discord_send_message",
		"Send a message to a Discord channel",
		SendMessageSchema(), SendMessage(bridge))

	builder.RegisterTool("discord_set_config",
		"Update Discord bot configuration (saved to ~/.orchestra/discord.json)",
		SetConfigSchema(), SetConfig(bridge))
}
