package common

import "testing"

func TestAnonymousUser(t *testing.T) {
	u := AnonymousUser()

	if u.ID != 1 {
		t.Errorf("expected ID 1, got %d", u.ID)
	}
	if u.Provider != "local" {
		t.Errorf("expected provider 'local', got %q", u.Provider)
	}
	if u.Name != "Skydiver" {
		t.Errorf("expected name 'Skydiver', got %q", u.Name)
	}
	if u.UnitSystem != "imperial" {
		t.Errorf("expected unit system 'imperial', got %q", u.UnitSystem)
	}
	if u.Locale != "en" {
		t.Errorf("expected locale 'en', got %q", u.Locale)
	}
}
