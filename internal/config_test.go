package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Enabled {
		t.Error("default config should be disabled")
	}
	if cfg.CommandPrefix != "!" {
		t.Errorf("default prefix should be !, got %q", cfg.CommandPrefix)
	}
	if len(cfg.AllowedUsers) != 0 {
		t.Error("default allowed users should be empty")
	}
}

func TestIsAllowed_EmptyList(t *testing.T) {
	cfg := DefaultConfig()
	if !cfg.IsAllowed("12345") {
		t.Error("empty allowed list should allow all users")
	}
}

func TestIsAllowed_WithList(t *testing.T) {
	cfg := &Config{AllowedUsers: []string{"111", "222", "333"}}
	if !cfg.IsAllowed("222") {
		t.Error("should allow user in list")
	}
	if cfg.IsAllowed("999") {
		t.Error("should deny user not in list")
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name  string
		cfg   Config
		valid bool
	}{
		{"empty", Config{}, false},
		{"webhook only", Config{WebhookURL: "https://example.com/webhook"}, true},
		{"bot token only", Config{BotToken: "token"}, false},
		{"bot+app", Config{BotToken: "token", ApplicationID: "app"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "discord.json")

	cfg := &Config{
		Enabled:       true,
		BotToken:      "test-token",
		ApplicationID: "12345",
		GuildID:       "67890",
		ChannelID:     "11111",
		CommandPrefix: "?",
		AllowedUsers:  []string{"user1", "user2"},
	}

	if err := cfg.SaveToFile(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded := LoadConfigFromFile(path)
	if loaded.BotToken != "test-token" {
		t.Errorf("BotToken = %q, want test-token", loaded.BotToken)
	}
	if loaded.CommandPrefix != "?" {
		t.Errorf("CommandPrefix = %q, want ?", loaded.CommandPrefix)
	}
	if len(loaded.AllowedUsers) != 2 {
		t.Errorf("AllowedUsers len = %d, want 2", len(loaded.AllowedUsers))
	}
	if !loaded.Enabled {
		t.Error("Enabled should be true")
	}
}

func TestLoadConfigFromFile_NotFound(t *testing.T) {
	cfg := LoadConfigFromFile("/nonexistent/path/discord.json")
	if cfg.Enabled {
		t.Error("should return default (disabled) when file not found")
	}
	if cfg.CommandPrefix != "!" {
		t.Errorf("should use default prefix, got %q", cfg.CommandPrefix)
	}
}

func TestSaveToFile_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "discord.json")

	cfg := &Config{BotToken: "test"}
	if err := cfg.SaveToFile(path); err != nil {
		t.Fatalf("Save should create dirs: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("File should exist: %v", err)
	}
}
