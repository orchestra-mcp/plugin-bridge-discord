package internal

import "testing"

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is longer than ten", 10, "this is..."},
		{"ab", 1, "a"},
		{"", 5, ""},
	}
	for _, tt := range tests {
		got := Truncate(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
	}
}

func TestSuccessEmbed(t *testing.T) {
	e := SuccessEmbed("Title", "Desc")
	if e.Color != 0x2ECC71 {
		t.Errorf("SuccessEmbed color = %x, want 2ECC71", e.Color)
	}
	if e.Title != "Title" {
		t.Errorf("Title = %q, want Title", e.Title)
	}
}

func TestErrorEmbed(t *testing.T) {
	e := ErrorEmbed("Err", "Something")
	if e.Color != 0xE74C3C {
		t.Errorf("ErrorEmbed color = %x, want E74C3C", e.Color)
	}
}

func TestInfoEmbed(t *testing.T) {
	e := InfoEmbed("Info", "Details")
	if e.Color != 0x3498DB {
		t.Errorf("InfoEmbed color = %x, want 3498DB", e.Color)
	}
}

func TestWarningEmbed(t *testing.T) {
	e := WarningEmbed("Warn", "Careful")
	if e.Color != 0xF39C12 {
		t.Errorf("WarningEmbed color = %x, want F39C12", e.Color)
	}
}

func TestToolEmbed(t *testing.T) {
	e := ToolEmbed("Read", "done", "file.go")
	if e.Color != 0x2ECC71 {
		t.Errorf("done ToolEmbed color = %x, want 2ECC71", e.Color)
	}

	e2 := ToolEmbed("Write", "error", "failed")
	if e2.Color != 0xE74C3C {
		t.Errorf("error ToolEmbed color = %x, want E74C3C", e2.Color)
	}
}

func TestPermissionEmbed(t *testing.T) {
	e := PermissionEmbed("Bash", "needs approval", "rm -rf")
	if e.Color != 0xF39C12 {
		t.Errorf("PermissionEmbed color = %x, want F39C12", e.Color)
	}
	if len(e.Fields) != 2 {
		t.Errorf("PermissionEmbed should have 2 fields, got %d", len(e.Fields))
	}
}

func TestPermissionResultEmbed(t *testing.T) {
	approved := PermissionResultEmbed("Approved", "req-1")
	if approved.Color != 0x2ECC71 {
		t.Errorf("Approved color = %x, want 2ECC71", approved.Color)
	}

	denied := PermissionResultEmbed("Denied", "req-2")
	if denied.Color != 0xE74C3C {
		t.Errorf("Denied color = %x, want E74C3C", denied.Color)
	}
}
