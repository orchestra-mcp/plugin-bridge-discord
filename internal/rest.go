package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const discordAPI = "https://discord.com/api/v10"

// RestClient communicates with the Discord REST API.
type RestClient struct {
	token  string
	client *http.Client
}

// NewRestClient creates a REST client for the Discord API.
func NewRestClient(token string) *RestClient {
	return &RestClient{
		token:  token,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

// SendMessage sends a message to a channel. Returns the message ID.
func (r *RestClient) SendMessage(channelID, content string, embed *Embed, components []Component) (string, error) {
	body := map[string]any{}
	if content != "" {
		body["content"] = content
	}
	if embed != nil {
		body["embeds"] = []Embed{*embed}
	}
	if len(components) > 0 {
		body["components"] = components
	}
	respBody, err := r.do("POST", fmt.Sprintf("/channels/%s/messages", channelID), body)
	if err != nil {
		return "", err
	}
	var msg struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(respBody, &msg)
	return msg.ID, nil
}

// EditMessage edits an existing message.
func (r *RestClient) EditMessage(channelID, messageID, content string, embed *Embed) error {
	body := map[string]any{}
	if content != "" {
		body["content"] = content
	}
	if embed != nil {
		body["embeds"] = []Embed{*embed}
	}
	_, err := r.do("PATCH", fmt.Sprintf("/channels/%s/messages/%s", channelID, messageID), body)
	return err
}

// RespondInteraction responds to a Discord interaction.
func (r *RestClient) RespondInteraction(id, token string, respType int, content string, embed *Embed) error {
	data := map[string]any{}
	if content != "" {
		data["content"] = content
	}
	if embed != nil {
		data["embeds"] = []Embed{*embed}
	}
	body := map[string]any{"type": respType, "data": data}
	_, err := r.do("POST", fmt.Sprintf("/interactions/%s/%s/callback", id, token), body)
	return err
}

// RegisterSlashCommands registers slash commands with Discord API.
func (r *RestClient) RegisterSlashCommands(appID, guildID string, cmds []SlashCommandDef) error {
	if guildID != "" {
		path := fmt.Sprintf("/applications/%s/guilds/%s/commands", appID, guildID)
		_, err := r.do("PUT", path, cmds)
		if err == nil {
			return nil
		}
		log.Printf("[discord] guild command registration failed, falling back to global: %v", err)
	}
	path := fmt.Sprintf("/applications/%s/commands", appID)
	_, err := r.do("PUT", path, cmds)
	return err
}

func (r *RestClient) do(method, path string, payload any) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshal payload: %w", err)
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, discordAPI+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bot "+r.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return respBody, fmt.Errorf("discord API %s %s: %d %s", method, path, resp.StatusCode, string(respBody))
	}
	return respBody, nil
}
