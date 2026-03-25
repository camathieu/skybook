package metadata

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/root-gg/skybook/common"
	"gorm.io/gorm"
)

// migrations returns the ordered list of database migrations.
// New migrations should be appended to this list — never modify existing entries.
func migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "202603241300_create_users",
			Migrate: func(tx *gorm.DB) error {
				if err := tx.AutoMigrate(&common.User{}); err != nil {
					return err
				}
				// Seed anonymous user for v1 single-user mode
				anon := common.AnonymousUser()
				return tx.Create(anon).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("users")
			},
		},
		{
			ID: "202603241301_create_jumps",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&common.Jump{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("jumps")
			},
		},
		{
			ID: "202603260001_add_packjob",
			Migrate: func(tx *gorm.DB) error {
				// AutoMigrate is idempotent — adds the column if missing, no-ops if present.
				return tx.AutoMigrate(&common.Jump{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropColumn(&common.Jump{}, "packjob")
			},
		},
	}
}
