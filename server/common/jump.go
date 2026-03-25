package common

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/datatypes"
)

// JumpType represents the discipline of the skydive.
type JumpType string

const (
	JumpTypeFF       JumpType = "FF"
	JumpTypeFS       JumpType = "FS"
	JumpTypeCRW      JumpType = "CRW"
	JumpTypeHOP      JumpType = "HOP"
	JumpTypeCF       JumpType = "CF"
	JumpTypeAFF      JumpType = "AFF"
	JumpTypeAFFI     JumpType = "AFFI"
	JumpTypeCamera   JumpType = "CAMERA"
	JumpTypeTandem   JumpType = "TANDEM"
	JumpTypeDemo     JumpType = "DEMO"
	JumpTypeXRW      JumpType = "XRW"
	JumpTypeAngle    JumpType = "ANGLE"
	JumpTypeTracking JumpType = "TRACKING"
	JumpTypeCP       JumpType = "CP"
	JumpTypeWingsuit JumpType = "WINGSUIT"
	JumpTypeOther    JumpType = "OTHER"
)

// IsValid returns true if the jump type is recognized.
func (jt JumpType) IsValid() bool {
	switch jt {
	case JumpTypeFF, JumpTypeFS, JumpTypeCRW, JumpTypeHOP,
		JumpTypeCF, JumpTypeAFF, JumpTypeAFFI, JumpTypeCamera,
		JumpTypeTandem, JumpTypeDemo, JumpTypeXRW,
		JumpTypeAngle, JumpTypeTracking, JumpTypeCP, JumpTypeWingsuit, JumpTypeOther:
		return true
	default:
		return false
	}
}

// AllJumpTypes returns all supported jump disciplines. Useful for exposing options to the frontend.
func AllJumpTypes() []JumpType {
	return []JumpType{
		JumpTypeFF, JumpTypeFS, JumpTypeCRW, JumpTypeHOP,
		JumpTypeCF, JumpTypeAFF, JumpTypeAFFI, JumpTypeCamera,
		JumpTypeTandem, JumpTypeDemo, JumpTypeXRW,
		JumpTypeAngle, JumpTypeTracking, JumpTypeCP, JumpTypeWingsuit, JumpTypeOther,
	}
}

// Jump represents a single skydive logbook entry.
type Jump struct {
	ID           uint                        `gorm:"primaryKey" json:"id"`
	UserID       uint                        `gorm:"not null;index;uniqueIndex:idx_user_number" json:"userId"`
	Number       uint                        `gorm:"not null;uniqueIndex:idx_user_number" json:"number"`
	Date         DateOnly                    `gorm:"not null;index" json:"date"`
	Dropzone     string                      `gorm:"size:255;not null;index" json:"dropzone"`
	Aircraft     string                      `gorm:"size:255" json:"aircraft,omitempty"`
	JumpType     JumpType                    `gorm:"size:32;not null;index" json:"jumpType"`
	Altitude     *uint                       `json:"altitude,omitempty"`
	FreefallTime *uint                       `json:"freefallTime,omitempty"`
	CanopySize   *uint                       `json:"canopySize,omitempty"`
	LO           string                      `gorm:"size:255" json:"lo,omitempty"`
	Event        string                      `gorm:"size:255" json:"event,omitempty"`
	Description  string                      `gorm:"type:text" json:"description,omitempty"`
	Links        datatypes.JSONSlice[string] `gorm:"type:text" json:"links,omitempty"`
	Landing      string                      `gorm:"size:32" json:"landing,omitempty"`
	NightJump    bool                        `gorm:"default:false" json:"nightJump"`
	OxygenJump   bool                        `gorm:"default:false" json:"oxygenJump"`
	CutAway      bool                        `gorm:"default:false" json:"cutaway"`
	Packjob      bool                        `gorm:"default:false" json:"packjob"`
	CreatedAt    time.Time                   `json:"createdAt"`
	UpdatedAt    time.Time                   `json:"updatedAt"`
}

// Validate checks that the required Jump fields are present and valid.
// Called by CreateJump and UpdateJump handlers before persisting.
func (j *Jump) Validate() error {
	if j.Date.IsZero() {
		return fmt.Errorf("date is required")
	}
	if strings.TrimSpace(j.Dropzone) == "" {
		return fmt.Errorf("dropzone is required")
	}
	if !j.JumpType.IsValid() {
		return fmt.Errorf("invalid jump_type")
	}
	return nil
}
