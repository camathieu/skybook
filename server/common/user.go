package common

import "time"

// User represents a SkyBook user.
// In v1 (anonymous mode), a single user (ID=1) is auto-created.
// In v6+ (multi-tenant), users are created via Google OAuth.
type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Provider   string    `gorm:"size:32;not null;uniqueIndex:idx_provider_id" json:"provider"`
	ProviderID string    `gorm:"size:255;uniqueIndex:idx_provider_id" json:"providerId"`
	Email      string    `gorm:"size:255" json:"email,omitempty"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Locale     string    `gorm:"size:10;not null;default:'en'" json:"locale"`
	UnitSystem string    `gorm:"size:10;not null;default:'imperial'" json:"unitSystem"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// AnonymousUser returns the default anonymous user for v1 single-user mode.
func AnonymousUser() *User {
	return &User{
		ID:         1,
		Provider:   "local",
		ProviderID: "",
		Name:       "Skydiver",
		Locale:     "en",
		UnitSystem: "imperial",
	}
}
