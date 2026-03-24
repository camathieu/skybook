package metadata

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/root-gg/skybook/common"
)

// CreateJump appends a new jump at the end of the user's logbook.
// The Number is automatically set to MAX(Number)+1.
func (b *Backend) CreateJump(jump *common.Jump) error {
	return b.db.Transaction(func(tx *gorm.DB) error {
		// Get next number
		var maxNumber uint
		tx.Model(&common.Jump{}).
			Where("user_id = ?", jump.UserID).
			Select("COALESCE(MAX(number), 0)").
			Scan(&maxNumber)

		jump.Number = maxNumber + 1
		return tx.Create(jump).Error
	})
}

// InsertJumpAt inserts a jump at a specific position in the user's logbook.
// All jumps with Number >= pos are shifted up by 1.
func (b *Backend) InsertJumpAt(jump *common.Jump, pos uint) error {
	if pos < 1 {
		return fmt.Errorf("position must be >= 1, got %d", pos)
	}

	return b.db.Transaction(func(tx *gorm.DB) error {
		// Check that pos is within bounds (allow appending at end+1)
		var count int64
		tx.Model(&common.Jump{}).Where("user_id = ?", jump.UserID).Count(&count)

		if pos > uint(count)+1 {
			return fmt.Errorf("position %d out of range (max %d)", pos, count+1)
		}

		// Shift existing jumps up. SQLite doesn't support UPDATE ... ORDER BY,
		// and we must update in descending order to avoid unique constraint violations.
		if count > 0 && pos <= uint(count) {
			var jumpsToShift []*common.Jump
			if err := tx.Where("user_id = ? AND number >= ?", jump.UserID, pos).
				Order("number DESC").
				Find(&jumpsToShift).Error; err != nil {
				return err
			}
			for _, j := range jumpsToShift {
				if err := tx.Model(j).Update("number", j.Number+1).Error; err != nil {
					return err
				}
			}
		}

		jump.Number = pos
		return tx.Create(jump).Error
	})
}

// DeleteJump removes a jump and renumbers subsequent jumps to maintain contiguity.
func (b *Backend) DeleteJump(jump *common.Jump) error {
	return b.db.Transaction(func(tx *gorm.DB) error {
		number := jump.Number
		userID := jump.UserID

		// Delete the jump
		if err := tx.Delete(jump).Error; err != nil {
			return err
		}

		// Shift subsequent jumps down. Must be ascending order to avoid unique constraint.
		var jumpsToShift []*common.Jump
		if err := tx.Where("user_id = ? AND number > ?", userID, number).
			Order("number ASC").
			Find(&jumpsToShift).Error; err != nil {
			return err
		}
		for _, j := range jumpsToShift {
			if err := tx.Model(j).Update("number", j.Number-1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetJump retrieves a single jump by ID for a given user.
func (b *Backend) GetJump(userID, jumpID uint) (*common.Jump, error) {
	jump := &common.Jump{}
	err := b.db.Where("id = ? AND user_id = ?", jumpID, userID).First(jump).Error
	if err != nil {
		return nil, err
	}
	return jump, nil
}

// GetJumpByNumber retrieves a single jump by number for a given user.
func (b *Backend) GetJumpByNumber(userID, number uint) (*common.Jump, error) {
	jump := &common.Jump{}
	err := b.db.Where("user_id = ? AND number = ?", userID, number).First(jump).Error
	if err != nil {
		return nil, err
	}
	return jump, nil
}

// UpdateJump updates a jump's fields (except Number — use InsertJumpAt/DeleteJump to reorder).
func (b *Backend) UpdateJump(jump *common.Jump) error {
	return b.db.Omit(clause.Associations).Save(jump).Error
}

// GetJumps retrieves a paginated list of jumps for a user.
func (b *Backend) GetJumps(userID uint, offset, limit int, sortBy, order string) ([]*common.Jump, int64, error) {
	var jumps []*common.Jump
	var total int64

	query := b.db.Where("user_id = ?", userID)

	// Count total matching
	query.Model(&common.Jump{}).Count(&total)

	// Apply sort
	if sortBy == "" {
		sortBy = "number"
	}
	if order == "" {
		order = "desc"
	}
	orderClause := fmt.Sprintf("%s %s", sortBy, order)

	err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&jumps).Error
	return jumps, total, err
}

// CountJumps returns the total number of jumps for a user.
func (b *Backend) CountJumps(userID uint) (int64, error) {
	var count int64
	err := b.db.Model(&common.Jump{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
