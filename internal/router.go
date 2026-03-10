package internal

import (
	"log"
	"strings"
)

// Router dispatches Discord events to matching handlers.
type Router struct {
	handlers       []Handler
	defaultHandler Handler
	prefix         string
}

// NewRouter creates a command router with the given prefix.
func NewRouter(prefix string) *Router {
	if prefix == "" {
		prefix = "!"
	}
	return &Router{prefix: prefix}
}

// Register adds a handler to the router.
func (r *Router) Register(h Handler) {
	r.handlers = append(r.handlers, h)
}

// SetDefault sets the fallback handler for unmatched prefix commands.
func (r *Router) SetDefault(h Handler) {
	r.defaultHandler = h
}

// RouteMessage dispatches a MESSAGE_CREATE event to the matching prefix handler.
func (r *Router) RouteMessage(msg MessageCreate, api HandlerAPI) {
	if msg.Author.Bot {
		return
	}
	if !api.Config().IsAllowed(msg.Author.ID) {
		return
	}
	content := strings.TrimSpace(msg.Content)
	if !strings.HasPrefix(content, r.prefix) {
		return
	}
	content = strings.TrimPrefix(content, r.prefix)

	for _, h := range r.handlers {
		if h.MatchesPrefix(content) {
			log.Printf("[discord] routing prefix command to %s", h.Name())
			go h.HandleMessage(&msg, api)
			return
		}
	}

	if r.defaultHandler != nil && content != "" {
		msg.Content = "chat " + content
		go r.defaultHandler.HandleMessage(&msg, api)
	}
}

// RouteInteraction dispatches an INTERACTION_CREATE event.
func (r *Router) RouteInteraction(ix InteractionCreate, api HandlerAPI) {
	if !api.Config().IsAllowed("") { // interactions don't always have user context easily
		return
	}
	switch ix.Type {
	case 2: // Slash command
		for _, h := range r.handlers {
			if h.MatchesSlash(ix.Data.Name) {
				go h.HandleInteraction(&ix, api)
				return
			}
		}
	case 3: // Button / component interaction
		for _, h := range r.handlers {
			if ih, ok := h.(InteractionHandler); ok && ih.MatchesCustomID(ix.Data.CustomID) {
				go ih.HandleComponentInteraction(&ix, api)
				return
			}
		}
	}
}

// SlashDefs collects slash command definitions from all handlers.
func (r *Router) SlashDefs() []SlashCommandDef {
	var defs []SlashCommandDef
	for _, h := range r.handlers {
		if def := h.SlashDef(); def != nil {
			defs = append(defs, *def)
		}
	}
	return defs
}
