package internal

import "testing"

func TestNewNotificationService_Nil(t *testing.T) {
	s := NewNotificationService(nil)
	if s != nil {
		t.Error("nil config should return nil service")
	}
}

func TestNewNotificationService_Disabled(t *testing.T) {
	s := NewNotificationService(&Config{Enabled: false})
	if s != nil {
		t.Error("disabled config should return nil service")
	}
}

func TestNewNotificationService_InvalidConfig(t *testing.T) {
	s := NewNotificationService(&Config{Enabled: true})
	if s != nil {
		t.Error("invalid config (no webhook or bot+channel) should return nil")
	}
}

func TestNewNotificationService_Webhook(t *testing.T) {
	s := NewNotificationService(&Config{
		Enabled:    true,
		WebhookURL: "https://example.com/webhook",
	})
	if s == nil {
		t.Error("webhook config should create service")
	}
}

func TestNewNotificationService_Bot(t *testing.T) {
	s := NewNotificationService(&Config{
		Enabled:   true,
		BotToken:  "token",
		ChannelID: "ch",
	})
	if s == nil {
		t.Error("bot+channel config should create service")
	}
}

func TestSendTransition_NilSafe(t *testing.T) {
	var s *NotificationService
	// Should not panic
	s.SendTransition("FEAT-1", "feature", "todo", "in-progress", "my-project")
}

func TestStatusEmoji(t *testing.T) {
	tests := map[string]string{
		"todo":        "📋",
		"in-progress": "🔨",
		"done":        "✅",
		"blocked":     "🚫",
		"unknown":     "🔄",
	}
	for status, want := range tests {
		got := statusEmoji(status)
		if got != want {
			t.Errorf("statusEmoji(%q) = %q, want %q", status, got, want)
		}
	}
}

func TestStatusColor(t *testing.T) {
	if c := statusColor("done"); c != 0x2ECC71 {
		t.Errorf("done color = %x, want 2ECC71", c)
	}
	if c := statusColor("blocked"); c != 0xE74C3C {
		t.Errorf("blocked color = %x, want E74C3C", c)
	}
	if c := statusColor("in-progress"); c != 0x3498DB {
		t.Errorf("in-progress color = %x, want 3498DB", c)
	}
}
