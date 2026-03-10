package internal

import (
	"encoding/json"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"
)

func TestExtractText_Nil(t *testing.T) {
	got := ExtractText(nil)
	if got != "" {
		t.Errorf("ExtractText(nil) = %q, want empty string", got)
	}
}

func TestExtractText_WithText(t *testing.T) {
	s, err := structpb.NewStruct(map[string]any{
		"text": "hello world",
	})
	if err != nil {
		t.Fatalf("NewStruct: %v", err)
	}
	got := ExtractText(s)
	if got != "hello world" {
		t.Errorf("ExtractText with text field = %q, want %q", got, "hello world")
	}
}

func TestExtractText_WithoutText(t *testing.T) {
	s, err := structpb.NewStruct(map[string]any{
		"status": "ok",
		"count":  float64(42),
	})
	if err != nil {
		t.Fatalf("NewStruct: %v", err)
	}
	got := ExtractText(s)

	// Should be valid JSON since there is no "text" field.
	var m map[string]any
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("ExtractText returned invalid JSON: %q, err: %v", got, err)
	}
	if m["status"] != "ok" {
		t.Errorf("JSON status = %v, want %q", m["status"], "ok")
	}
	if m["count"] != float64(42) {
		t.Errorf("JSON count = %v, want 42", m["count"])
	}
}

func TestExtractText_EmptyStruct(t *testing.T) {
	s, err := structpb.NewStruct(map[string]any{})
	if err != nil {
		t.Fatalf("NewStruct: %v", err)
	}
	got := ExtractText(s)
	if got != "{}" {
		t.Errorf("ExtractText(empty struct) = %q, want %q", got, "{}")
	}
}

func TestExtractText_NestedStruct(t *testing.T) {
	s, err := structpb.NewStruct(map[string]any{
		"data": map[string]any{
			"name": "test",
			"tags": []any{"a", "b"},
		},
		"active": true,
	})
	if err != nil {
		t.Fatalf("NewStruct: %v", err)
	}
	got := ExtractText(s)

	// Should be valid JSON containing nested data.
	var m map[string]any
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("ExtractText returned invalid JSON: %q, err: %v", got, err)
	}
	data, ok := m["data"].(map[string]any)
	if !ok {
		t.Fatalf("data field is not a map: %T", m["data"])
	}
	if data["name"] != "test" {
		t.Errorf("nested name = %v, want %q", data["name"], "test")
	}
	tags, ok := data["tags"].([]any)
	if !ok {
		t.Fatalf("tags is not a slice: %T", data["tags"])
	}
	if len(tags) != 2 {
		t.Errorf("tags length = %d, want 2", len(tags))
	}
	if m["active"] != true {
		t.Errorf("active = %v, want true", m["active"])
	}
}
