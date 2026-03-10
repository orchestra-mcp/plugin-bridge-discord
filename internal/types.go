package internal

import "encoding/json"

// MessageCreate represents a Discord MESSAGE_CREATE gateway event.
type MessageCreate struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
	Author    Author `json:"author"`
}

// Author represents a Discord message author.
type Author struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bot      bool   `json:"bot"`
}

// InteractionCreate represents a Discord INTERACTION_CREATE gateway event.
type InteractionCreate struct {
	ID      string          `json:"id"`
	Type    int             `json:"type"` // 2=slash command, 3=message component (button)
	Token   string          `json:"token"`
	Data    InteractionData `json:"data"`
	Message *MessageRef     `json:"message,omitempty"`
}

// InteractionData holds the command data from an interaction.
type InteractionData struct {
	Name     string              `json:"name"`
	CustomID string              `json:"custom_id"`
	Options  []InteractionOption `json:"options"`
}

// InteractionOption is a single option in a slash command.
type InteractionOption struct {
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
}

// OptionString extracts a string value from the option.
func (o InteractionOption) OptionString() string {
	var s string
	_ = json.Unmarshal(o.Value, &s)
	return s
}

// OptionBool extracts a bool value from the option.
func (o InteractionOption) OptionBool() bool {
	var b bool
	_ = json.Unmarshal(o.Value, &b)
	return b
}

// MessageRef is a reference to a message in a channel.
type MessageRef struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
}

// SlashCommandDef defines a slash command for registration with Discord API.
type SlashCommandDef struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Options     []SlashOptionDef `json:"options,omitempty"`
}

// SlashOptionDef defines an option for a slash command.
type SlashOptionDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"` // 3=STRING, 5=BOOL
	Required    bool   `json:"required"`
}

// Component represents a Discord message component (buttons, action rows).
type Component struct {
	Type       int         `json:"type"`                 // 1=ActionRow, 2=Button
	Components []Component `json:"components,omitempty"` // children for ActionRow
	Style      int         `json:"style,omitempty"`      // 1=Primary, 3=Success, 4=Danger
	Label      string      `json:"label,omitempty"`
	CustomID   string      `json:"custom_id,omitempty"`
}

// ActionRow creates a component action row containing the given buttons.
func ActionRow(buttons ...Component) Component {
	return Component{Type: 1, Components: buttons}
}

// Button creates a button component.
func Button(label, customID string, style int) Component {
	return Component{Type: 2, Style: style, Label: label, CustomID: customID}
}

// Button styles
const (
	ButtonPrimary = 1
	ButtonSuccess = 3
	ButtonDanger  = 4
)

// Interaction response types
const (
	InteractionResponseMessage  = 4 // respond with message
	InteractionResponseDeferred = 5 // ACK, edit later
	InteractionResponseUpdate   = 7 // update existing message
)

// Embed represents a Discord embed.
type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Color       int          `json:"color"`
	Fields      []EmbedField `json:"fields,omitempty"`
	Footer      *EmbedFooter `json:"footer,omitempty"`
	Timestamp   string       `json:"timestamp,omitempty"`
}

// EmbedField represents a field in a Discord embed.
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// EmbedFooter represents footer text in a Discord embed.
type EmbedFooter struct {
	Text string `json:"text"`
}
