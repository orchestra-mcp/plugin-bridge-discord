package internal

import (
	"encoding/json"
	"testing"
)

func TestOptionString(t *testing.T) {
	raw, _ := json.Marshal("hello")
	opt := InteractionOption{Name: "test", Value: raw}
	if got := opt.OptionString(); got != "hello" {
		t.Errorf("OptionString() = %q, want hello", got)
	}
}

func TestOptionBool(t *testing.T) {
	raw, _ := json.Marshal(true)
	opt := InteractionOption{Name: "flag", Value: raw}
	if !opt.OptionBool() {
		t.Error("OptionBool() should be true")
	}
}

func TestActionRow(t *testing.T) {
	row := ActionRow(
		Button("OK", "btn_ok", ButtonSuccess),
		Button("Cancel", "btn_cancel", ButtonDanger),
	)
	if row.Type != 1 {
		t.Errorf("ActionRow type should be 1, got %d", row.Type)
	}
	if len(row.Components) != 2 {
		t.Errorf("ActionRow should have 2 components, got %d", len(row.Components))
	}
}

func TestButton(t *testing.T) {
	btn := Button("Click", "custom_id", ButtonPrimary)
	if btn.Type != 2 {
		t.Errorf("Button type should be 2, got %d", btn.Type)
	}
	if btn.Label != "Click" {
		t.Errorf("Label = %q, want Click", btn.Label)
	}
	if btn.Style != ButtonPrimary {
		t.Errorf("Style = %d, want %d", btn.Style, ButtonPrimary)
	}
}

func TestEmbedJSON(t *testing.T) {
	embed := Embed{
		Title:       "Test",
		Description: "Desc",
		Color:       0x2ECC71,
		Fields:      []EmbedField{{Name: "F1", Value: "V1", Inline: true}},
		Footer:      &EmbedFooter{Text: "Footer"},
	}
	data, err := json.Marshal(embed)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Embed
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if decoded.Title != "Test" {
		t.Errorf("Title = %q, want Test", decoded.Title)
	}
	if len(decoded.Fields) != 1 {
		t.Fatalf("Fields len = %d, want 1", len(decoded.Fields))
	}
	if decoded.Fields[0].Name != "F1" {
		t.Errorf("Field name = %q, want F1", decoded.Fields[0].Name)
	}
}
