package metadata

import "github.com/go-gormigrate/gormigrate/v2"

// migrations returns the ordered list of database migrations.
// New migrations should be appended to this list.
func migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		// Migrations will be added here as models are created.
		// Example:
		// {
		//     ID: "202603241200_create_jumps",
		//     Migrate: func(tx *gorm.DB) error {
		//         return tx.AutoMigrate(&common.Jump{})
		//     },
		//     Rollback: func(tx *gorm.DB) error {
		//         return tx.Migrator().DropTable("jumps")
		//     },
		// },
	}
}
