package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds Discord bot configuration.
type Config struct {
	Enabled       bool     `json:"enabled"`
	BotToken      string   `json:"bot_token"`
	ClientID      string   `json:"client_id"`
	ClientSecret  string   `json:"client_secret"`
	ApplicationID string   `json:"application_id"`
	GuildID       string   `json:"guild_id"`
	ChannelID     string   `json:"channel_id"`
	CommandPrefix string   `json:"command_prefix"`
	WebhookURL    string   `json:"webhook_url"`
	AllowedUsers  []string `json:"allowed_users"`
}

// DefaultConfig returns default Discord configuration.
func DefaultConfig() *Config {
	return &Config{
		Enabled:       false,
		CommandPrefix: "!",
		AllowedUsers:  []string{},
	}
}

// IsAllowed checks if a Discord user ID is in the allowed list.
// If AllowedUsers is empty, all users are allowed.
func (c *Config) IsAllowed(userID string) bool {
	if len(c.AllowedUsers) == 0 {
		return true
	}
	for _, id := range c.AllowedUsers {
		if id == userID {
			return true
		}
	}
	return false
}

// IsValid checks if minimum required fields are present.
func (c *Config) IsValid() bool {
	if c.BotToken != "" && c.ApplicationID != "" {
		return true
	}
	if c.WebhookURL != "" {
		return true
	}
	return false
}

// ConfigPath returns the default config file path.
func ConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".orchestra", "discord.json")
}

// LoadConfig loads config from the default path.
func LoadConfig() *Config {
	return LoadConfigFromFile(ConfigPath())
}

// LoadConfigFromFile loads config from a specific path.
func LoadConfigFromFile(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig()
	}
	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return DefaultConfig()
	}
	if cfg.CommandPrefix == "" {
		cfg.CommandPrefix = "!"
	}
	return cfg
}

// Save writes config to the default path.
func (c *Config) Save() error {
	return c.SaveToFile(ConfigPath())
}

// SaveToFile writes config to a specific path.
func (c *Config) SaveToFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
