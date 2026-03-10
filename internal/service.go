package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NotificationService sends Discord notifications for workflow transitions.
type NotificationService struct {
	config *Config
	client *http.Client
}

// NewNotificationService creates a Discord notification service.
func NewNotificationService(cfg *Config) *NotificationService {
	if cfg == nil || !cfg.Enabled {
		return nil
	}
	if cfg.WebhookURL == "" && (cfg.BotToken == "" || cfg.ChannelID == "") {
		return nil
	}
	return &NotificationService{
		config: cfg,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// SendTransition sends a workflow transition notification to Discord.
func (s *NotificationService) SendTransition(taskID, taskType, from, to, project string) {
	if s == nil {
		return
	}
	emoji := statusEmoji(to)
	color := statusColor(to)
	title := fmt.Sprintf("%s %s -> %s", emoji, taskID, to)
	desc := fmt.Sprintf("%s **%s** `%s` -> `%s`", emoji, taskID, from, to)
	fields := []EmbedField{
		{Name: "Task", Value: fmt.Sprintf("`%s` %s", taskID, taskType), Inline: true},
		{Name: "Status", Value: fmt.Sprintf("`%s` -> `%s`", from, to), Inline: true},
	}
	embed := Embed{
		Title:       title,
		Description: desc,
		Color:       color,
		Fields:      fields,
		Footer:      &EmbedFooter{Text: fmt.Sprintf("Project: %s", project)},
	}
	go s.send(embed)
}

func (s *NotificationService) send(embed Embed) {
	body := map[string]any{"embeds": []Embed{embed}}
	data, err := json.Marshal(body)
	if err != nil {
		return
	}
	if s.config.WebhookURL != "" {
		req, _ := http.NewRequest("POST", s.config.WebhookURL, bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := s.client.Do(req)
		if resp != nil {
			resp.Body.Close()
		}
		return
	}
	url := fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", s.config.ChannelID)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(data))
	req.Header.Set("Authorization", "Bot "+s.config.BotToken)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := s.client.Do(req)
	if resp != nil {
		resp.Body.Close()
	}
}

func statusEmoji(status string) string {
	switch status {
	case "todo":
		return "\xf0\x9f\x93\x8b"
	case "in-progress":
		return "\xf0\x9f\x94\xa8"
	case "in-testing":
		return "\xf0\x9f\x94\xac"
	case "in-docs":
		return "\xe2\x9c\x8d\xef\xb8\x8f"
	case "in-review":
		return "\xf0\x9f\x91\x80"
	case "done":
		return "\xe2\x9c\x85"
	case "blocked":
		return "\xf0\x9f\x9a\xab"
	case "rejected", "needs-edits":
		return "\xe2\x9d\x8c"
	default:
		return "\xf0\x9f\x94\x84"
	}
}

func statusColor(status string) int {
	switch status {
	case "todo":
		return 0x546E7A
	case "in-progress":
		return 0x3498DB
	case "in-testing":
		return 0xF39C12
	case "in-docs":
		return 0x9B59B6
	case "in-review":
		return 0xE67E22
	case "done":
		return 0x2ECC71
	case "blocked", "rejected", "needs-edits":
		return 0xE74C3C
	default:
		return 0x7289DA
	}
}
