package internal

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// Discord limits
const (
	MaxEmbedDesc   = 4096
	MaxFieldValue  = 1024
	MaxMessageText = 2000
	// Safe limits (leave room for formatting)
	SafeEmbedDesc  = 3900
	SafeFieldValue = 950
	SafeMessageLen = 1990
)

// SuccessEmbed creates a green embed for successful operations.
func SuccessEmbed(title, desc string) *Embed {
	return &Embed{Title: Truncate(title, 256), Description: Truncate(desc, SafeEmbedDesc), Color: 0x2ECC71}
}

// ErrorEmbed creates a red embed for errors.
func ErrorEmbed(title, desc string) *Embed {
	return &Embed{Title: Truncate(title, 256), Description: Truncate(desc, SafeEmbedDesc), Color: 0xE74C3C}
}

// InfoEmbed creates a blue embed for informational messages.
func InfoEmbed(title, desc string) *Embed {
	return &Embed{Title: Truncate(title, 256), Description: Truncate(desc, SafeEmbedDesc), Color: 0x3498DB}
}

// WarningEmbed creates an orange embed for warnings.
func WarningEmbed(title, desc string) *Embed {
	return &Embed{Title: Truncate(title, 256), Description: Truncate(desc, SafeEmbedDesc), Color: 0xF39C12}
}

// ToolEmbed creates a colored embed for tool execution status.
// Accepts raw tool name + detail, or parses JSON action payload.
func ToolEmbed(tool, status, detail string) *Embed {
	color := 0x3498DB
	emoji := "\xf0\x9f\x94\xa7"
	switch status {
	case "done":
		color = 0x2ECC71
		emoji = "\xe2\x9c\x85"
	case "error":
		color = 0xE74C3C
		emoji = "\xe2\x9d\x8c"
	}
	return &Embed{
		Title:       Truncate(fmt.Sprintf("%s %s", emoji, tool), 256),
		Description: Truncate(detail, SafeEmbedDesc),
		Color:       color,
	}
}

// ActionEmbed parses a raw action JSON payload and returns a human-readable embed.
func ActionEmbed(rawJSON string) *Embed {
	var action struct {
		Tool      string `json:"tool"`
		Input     string `json:"input"`
		Status    string `json:"status"`
		ToolUseID string `json:"toolUseId"`
		Result    string `json:"result"`
	}
	if err := json.Unmarshal([]byte(rawJSON), &action); err != nil {
		return ToolEmbed("Tool", "running", Truncate(rawJSON, SafeEmbedDesc))
	}

	emoji := toolEmoji(action.Tool)
	color := 0x3498DB // blue = running
	if action.Status == "done" {
		color = 0x2ECC71 // green = done
	}

	title := fmt.Sprintf("%s %s", emoji, humanToolName(action.Tool))
	desc := formatToolInput(action.Tool, action.Input)

	embed := &Embed{
		Title:       Truncate(title, 256),
		Description: Truncate(desc, SafeEmbedDesc),
		Color:       color,
	}

	if action.Result != "" && action.Status == "done" {
		embed.Fields = []EmbedField{
			{Name: "Result", Value: Truncate(formatCodeBlock(action.Result), SafeFieldValue)},
		}
	}

	return embed
}

// PermissionEmbed creates an embed for permission requests.
func PermissionEmbed(toolName, reason, input string) *Embed {
	return &Embed{
		Title:       Truncate(fmt.Sprintf("\xf0\x9f\x94\x90 Permission: %s", toolName), 256),
		Description: Truncate(reason, SafeEmbedDesc),
		Color:       0xF39C12,
		Fields: []EmbedField{
			{Name: "Tool", Value: fmt.Sprintf("`%s`", toolName), Inline: true},
			{Name: "Input", Value: Truncate(formatCodeBlock(input), SafeFieldValue)},
		},
	}
}

// PermissionResultEmbed creates an embed showing the result of a permission decision.
func PermissionResultEmbed(status, requestID string) *Embed {
	color := 0x2ECC71
	emoji := "\xe2\x9c\x85"
	if status == "Denied" {
		color = 0xE74C3C
		emoji = "\xe2\x9d\x8c"
	}
	return &Embed{
		Title: fmt.Sprintf("%s Permission %s", emoji, status),
		Color: color,
	}
}

// Truncate limits a string to maxLen characters with ellipsis.
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen < 4 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// --- Internal helpers ---

func toolEmoji(tool string) string {
	switch tool {
	case "Read":
		return "\xf0\x9f\x93\x96"
	case "Write":
		return "\xf0\x9f\x93\x9d"
	case "Edit":
		return "\xe2\x9c\x8f\xef\xb8\x8f"
	case "Bash":
		return "\xf0\x9f\x92\xbb"
	case "Grep":
		return "\xf0\x9f\x94\x8d"
	case "Glob":
		return "\xf0\x9f\x93\x82"
	case "Task":
		return "\xf0\x9f\xa4\x96"
	case "TodoWrite":
		return "\xf0\x9f\x93\x8b"
	case "WebFetch":
		return "\xf0\x9f\x8c\x90"
	case "WebSearch":
		return "\xf0\x9f\x94\x8e"
	default:
		if strings.HasPrefix(tool, "mcp__") {
			return "\xf0\x9f\x94\x8c"
		}
		return "\xf0\x9f\x94\xa7"
	}
}

func humanToolName(tool string) string {
	// MCP tools: mcp__server__tool_name -> server/tool_name
	if strings.HasPrefix(tool, "mcp__") {
		parts := strings.SplitN(tool, "__", 3)
		if len(parts) == 3 {
			return parts[1] + "/" + strings.ReplaceAll(parts[2], "_", " ")
		}
	}
	switch tool {
	case "Read":
		return "Reading file"
	case "Write":
		return "Writing file"
	case "Edit":
		return "Editing file"
	case "Bash":
		return "Running command"
	case "Grep":
		return "Searching code"
	case "Glob":
		return "Finding files"
	case "Task":
		return "Sub-agent task"
	case "TodoWrite":
		return "Updating todo list"
	case "WebFetch":
		return "Fetching URL"
	case "WebSearch":
		return "Web search"
	default:
		return tool
	}
}

func formatToolInput(tool, rawInput string) string {
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(rawInput), &input); err != nil {
		return Truncate(rawInput, 200)
	}

	switch tool {
	case "Read":
		fp, _ := input["file_path"].(string)
		if fp != "" {
			return fmt.Sprintf("`%s`", shortPath(fp))
		}
	case "Write":
		fp, _ := input["file_path"].(string)
		if fp != "" {
			return fmt.Sprintf("`%s`", shortPath(fp))
		}
	case "Edit":
		fp, _ := input["file_path"].(string)
		old, _ := input["old_string"].(string)
		if fp != "" {
			desc := fmt.Sprintf("`%s`", shortPath(fp))
			if old != "" {
				desc += fmt.Sprintf("\n```\n%s\n```", Truncate(old, 100))
			}
			return desc
		}
	case "Bash":
		cmd, _ := input["command"].(string)
		if cmd != "" {
			return fmt.Sprintf("```\n%s\n```", Truncate(cmd, 300))
		}
	case "Grep":
		pat, _ := input["pattern"].(string)
		path, _ := input["path"].(string)
		desc := fmt.Sprintf("Pattern: `%s`", pat)
		if path != "" {
			desc += fmt.Sprintf(" in `%s`", shortPath(path))
		}
		return desc
	case "Glob":
		pat, _ := input["pattern"].(string)
		return fmt.Sprintf("Pattern: `%s`", pat)
	case "Task":
		desc, _ := input["description"].(string)
		agentType, _ := input["subagent_type"].(string)
		if desc != "" {
			s := desc
			if agentType != "" {
				s = fmt.Sprintf("[%s] %s", agentType, desc)
			}
			return Truncate(s, 300)
		}
	case "TodoWrite":
		return "Updating task list"
	case "WebSearch":
		q, _ := input["query"].(string)
		if q != "" {
			return fmt.Sprintf("Query: `%s`", Truncate(q, 200))
		}
	case "WebFetch":
		u, _ := input["url"].(string)
		if u != "" {
			return Truncate(u, 300)
		}
	default:
		// MCP tools -- show key args compactly
		if strings.HasPrefix(tool, "mcp__") {
			return formatMCPInput(input)
		}
	}

	return ""
}

func formatMCPInput(input map[string]interface{}) string {
	var parts []string
	for k, v := range input {
		s := fmt.Sprintf("%v", v)
		if len(s) > 80 {
			s = s[:77] + "..."
		}
		parts = append(parts, fmt.Sprintf("**%s:** %s", k, s))
	}
	return Truncate(strings.Join(parts, "\n"), 500)
}

func shortPath(fp string) string {
	// Show just filename or last 2 path components
	base := filepath.Base(fp)
	dir := filepath.Base(filepath.Dir(fp))
	if dir == "." || dir == "/" {
		return base
	}
	return dir + "/" + base
}

func formatCodeBlock(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if len(s) > SafeFieldValue-10 {
		s = s[:SafeFieldValue-15] + "..."
	}
	return "```\n" + s + "\n```"
}
