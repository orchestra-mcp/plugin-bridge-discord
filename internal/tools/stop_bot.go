package tools

import (
	"context"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// StopBotSchema returns the JSON schema for the stop_discord_bot tool.
func StopBotSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	})
	return s
}

// StopBot returns a tool handler that stops the Discord bot.
func StopBot(bridge *DiscordBridge) func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if bridge.Plugin.Bot == nil {
			return helpers.ErrorResult("bot_error", "Bot not initialized"), nil
		}
		if !bridge.Plugin.Bot.IsRunning() {
			return helpers.TextResult("Discord bot is not running"), nil
		}
		bridge.Plugin.Bot.Stop()
		return helpers.TextResult("Discord bot stopped"), nil
	}
}
