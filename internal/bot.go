package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

// ToolCaller is a function that invokes MCP tools by name.
type ToolCaller func(name string, args map[string]any) (string, error)

// HandlerRegistrar is a callback that registers handlers on a router.
// This allows external packages to wire up handlers without import cycles.
type HandlerRegistrar func(r *Router)

// Bot manages Discord integration -- gateway, router, handlers.
type Bot struct {
	config    *Config
	gateway   *Gateway
	router    *Router
	rest      *RestClient
	service   *NotificationService
	caller    ToolCaller
	registrar HandlerRegistrar
	mu        sync.RWMutex
	running   bool
}

// NewBot creates a Discord bot from configuration.
func NewBot(cfg *Config, caller ToolCaller) *Bot {
	return &Bot{
		config:  cfg,
		service: NewNotificationService(cfg),
		caller:  caller,
	}
}

// SetHandlerRegistrar sets the callback used to register handlers on Start.
func (b *Bot) SetHandlerRegistrar(fn HandlerRegistrar) {
	b.registrar = fn
}

// Start connects to Discord, registers handlers, and starts routing events.
func (b *Bot) Start(ctx context.Context) error {
	if b.config == nil || !b.config.Enabled || b.config.BotToken == "" {
		log.Println("[discord] bot disabled or not configured")
		return nil
	}

	gw, err := ConnectGateway(b.config.BotToken)
	if err != nil {
		return fmt.Errorf("connect gateway: %w", err)
	}
	b.gateway = gw
	b.rest = NewRestClient(b.config.BotToken)

	prefix := b.config.CommandPrefix
	if prefix == "" {
		prefix = "!"
	}
	b.router = NewRouter(prefix)

	// Register handlers via external registrar
	if b.registrar != nil {
		b.registrar(b.router)
	}

	// Register slash commands
	if b.config.ApplicationID != "" {
		defs := b.router.SlashDefs()
		if len(defs) > 0 {
			if err := b.rest.RegisterSlashCommands(b.config.ApplicationID, b.config.GuildID, defs); err != nil {
				log.Printf("[discord] failed to register slash commands: %v", err)
			}
		}
	}

	gw.SetEventHandler(b.onGatewayEvent)
	b.running = true

	log.Printf("[discord] bot started (%d handlers)", len(b.router.handlers))

	// Wait for context cancellation
	<-ctx.Done()
	b.Stop()
	return nil
}

// Stop gracefully stops the Discord bot.
func (b *Bot) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.gateway != nil {
		b.gateway.Close()
		b.gateway = nil
	}
	b.running = false
	log.Println("[discord] bot stopped")
}

// IsRunning returns whether the bot is active.
func (b *Bot) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.running
}

func (b *Bot) onGatewayEvent(eventType string, data json.RawMessage) {
	switch eventType {
	case "MESSAGE_CREATE":
		var msg MessageCreate
		if err := json.Unmarshal(data, &msg); err != nil {
			return
		}
		b.router.RouteMessage(msg, b)
	case "INTERACTION_CREATE":
		var ix InteractionCreate
		if err := json.Unmarshal(data, &ix); err != nil {
			return
		}
		b.router.RouteInteraction(ix, b)
	}
}

// --- HandlerAPI implementation ---

// SendToChannel implements HandlerAPI.
func (b *Bot) SendToChannel(channelID, content string, embed *Embed) error {
	if b.rest == nil {
		return fmt.Errorf("REST client not initialized")
	}
	_, err := b.rest.SendMessage(channelID, content, embed, nil)
	return err
}

// SendComponents implements HandlerAPI.
func (b *Bot) SendComponents(channelID, content string, embed *Embed, components []Component) (string, error) {
	if b.rest == nil {
		return "", fmt.Errorf("REST client not initialized")
	}
	return b.rest.SendMessage(channelID, content, embed, components)
}

// EditMessage implements HandlerAPI.
func (b *Bot) EditMessage(channelID, messageID, content string, embed *Embed) error {
	if b.rest == nil {
		return fmt.Errorf("REST client not initialized")
	}
	return b.rest.EditMessage(channelID, messageID, content, embed)
}

// RespondInteraction implements HandlerAPI.
func (b *Bot) RespondInteraction(id, token string, respType int, content string, embed *Embed) error {
	if b.rest == nil {
		return fmt.Errorf("REST client not initialized")
	}
	return b.rest.RespondInteraction(id, token, respType, content, embed)
}

// ChannelID implements HandlerAPI.
func (b *Bot) ChannelID() string { return b.config.ChannelID }

// Config implements HandlerAPI.
func (b *Bot) Config() *Config { return b.config }

// CallTool implements HandlerAPI -- invokes MCP tools via cross-plugin calls.
func (b *Bot) CallTool(name string, args map[string]any) (string, error) {
	if b.caller == nil {
		return "", fmt.Errorf("tool caller not configured")
	}
	return b.caller(name, args)
}
