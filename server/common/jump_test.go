package common

import (
	"encoding/json"
	"testing"
)

func TestIsValidJumpType(t *testing.T) {
	valid := []string{"FF", "WS", "FS", "CRW", "HOP", "CF", "AFF", "TANDEM", "DEMO", "XRW", "ANGLE", "TRACKING", "CP", "WINGSUIT", "OTHER"}
	for _, jt := range valid {
		if !JumpType(jt).IsValid() {
			t.Errorf("expected %q to be valid", jt)
		}
	}

	invalid := []string{"", "ff", "INVALID", "belly", "freefly"}
	for _, jt := range invalid {
		if JumpType(jt).IsValid() {
			t.Errorf("expected %q to be invalid", jt)
		}
	}
}

func TestJumpLinksJSON(t *testing.T) {
	jump := &Jump{
		Links: []string{"https://youtube.com/watch?v=abc", "https://example.com/photo.jpg"},
	}

	data, err := json.Marshal(jump)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Jump
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(decoded.Links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(decoded.Links))
	}
	if decoded.Links[0] != "https://youtube.com/watch?v=abc" {
		t.Errorf("link[0] = %q", decoded.Links[0])
	}
}

func TestJumpNullableFields(t *testing.T) {
	// Zero value (nil pointer) vs explicit zero
	jump := &Jump{}

	if jump.Altitude != nil {
		t.Error("expected nil Altitude for zero value")
	}

	alt := uint(0)
	jump.Altitude = &alt
	if jump.Altitude == nil || *jump.Altitude != 0 {
		t.Error("expected *Altitude == 0")
	}

	data, err := json.Marshal(jump)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// When nil, altitude should be omitted from JSON
	jump2 := &Jump{}
	data2, _ := json.Marshal(jump2)
	str := string(data2)
	if containsField(str, "altitude") {
		t.Error("nil altitude should be omitted from JSON")
	}

	// When set (even to 0), it should be present
	str = string(data)
	if !containsField(str, "altitude") {
		t.Error("non-nil altitude should be present in JSON")
	}
}

func containsField(jsonStr, field string) bool {
	return len(jsonStr) > 0 && json.Valid([]byte(jsonStr)) && contains(jsonStr, `"`+field+`"`)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
