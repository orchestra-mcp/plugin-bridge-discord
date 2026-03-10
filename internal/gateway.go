package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Gateway maintains a WebSocket connection to Discord.
type Gateway struct {
	token     string
	conn      *websocket.Conn
	heartbeat time.Duration
	seq       *int64
	done      chan struct{}
	mu        sync.Mutex
	onEvent   EventCallback
}

// EventCallback is called by the gateway for DISPATCH events.
type EventCallback func(eventType string, data json.RawMessage)

// SetEventHandler sets the callback for DISPATCH (op 0) gateway events.
func (g *Gateway) SetEventHandler(cb EventCallback) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.onEvent = cb
}

type gatewayPayload struct {
	Op int             `json:"op"`
	D  json.RawMessage `json:"d,omitempty"`
	S  *int64          `json:"s,omitempty"`
	T  string          `json:"t,omitempty"`
}

type helloData struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

// ConnectGateway connects to Discord Gateway and keeps the bot online.
func ConnectGateway(botToken string) (*Gateway, error) {
	gatewayURL, err := getGatewayURL()
	if err != nil {
		return nil, fmt.Errorf("get gateway URL: %w", err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(gatewayURL+"/?v=10&encoding=json", nil)
	if err != nil {
		return nil, fmt.Errorf("dial gateway: %w", err)
	}

	g := &Gateway{
		token: botToken,
		conn:  conn,
		done:  make(chan struct{}),
	}

	if err := g.readHello(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("read hello: %w", err)
	}

	if err := g.sendIdentify(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("identify: %w", err)
	}

	go g.heartbeatLoop()
	go g.readLoop()

	return g, nil
}

// Close disconnects from the gateway.
func (g *Gateway) Close() {
	select {
	case <-g.done:
		return
	default:
		close(g.done)
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.conn != nil {
		g.conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		g.conn.Close()
	}
}

func getGatewayURL() (string, error) {
	resp, err := http.Get("https://discord.com/api/v10/gateway")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	return data.URL, nil
}

func (g *Gateway) readHello() error {
	var payload gatewayPayload
	if err := g.conn.ReadJSON(&payload); err != nil {
		return err
	}
	if payload.Op != 10 {
		return fmt.Errorf("expected op 10 (Hello), got op %d", payload.Op)
	}
	var hello helloData
	if err := json.Unmarshal(payload.D, &hello); err != nil {
		return err
	}
	g.heartbeat = time.Duration(hello.HeartbeatInterval) * time.Millisecond
	return nil
}

func (g *Gateway) sendIdentify() error {
	identify := map[string]any{
		"token":   g.token,
		"intents": (1 << 0) | (1 << 9) | (1 << 15),
		"properties": map[string]string{
			"os":      "linux",
			"browser": "orchestra",
			"device":  "orchestra",
		},
		"presence": map[string]any{
			"status": "online",
			"activities": []map[string]any{
				{"name": "Orchestra MCP", "type": 0},
			},
		},
	}
	data, _ := json.Marshal(identify)
	return g.conn.WriteJSON(gatewayPayload{Op: 2, D: json.RawMessage(data)})
}

func (g *Gateway) heartbeatLoop() {
	ticker := time.NewTicker(g.heartbeat)
	defer ticker.Stop()
	for {
		select {
		case <-g.done:
			return
		case <-ticker.C:
			g.mu.Lock()
			var seqData json.RawMessage
			if g.seq != nil {
				seqData, _ = json.Marshal(*g.seq)
			} else {
				seqData = json.RawMessage("null")
			}
			err := g.conn.WriteJSON(gatewayPayload{Op: 1, D: seqData})
			g.mu.Unlock()
			if err != nil {
				return
			}
		}
	}
}

func (g *Gateway) readLoop() {
	for {
		select {
		case <-g.done:
			return
		default:
		}
		var payload gatewayPayload
		if err := g.conn.ReadJSON(&payload); err != nil {
			return
		}
		if payload.S != nil {
			g.mu.Lock()
			g.seq = payload.S
			g.mu.Unlock()
		}
		if payload.Op == 0 && g.onEvent != nil {
			g.onEvent(payload.T, payload.D)
		}
	}
}
