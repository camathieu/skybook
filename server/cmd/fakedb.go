package cmd

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/root-gg/skybook/common"
	"github.com/root-gg/skybook/metadata"
	"github.com/spf13/cobra"
)

var (
	fakedbJumps  int
	fakedbOutput string
)

var fakedbCmd = &cobra.Command{
	Use:   "fakedb",
	Short: "Generate randomized test data",
	Long:  "Creates a SQLite database populated with realistic chronological jump data for testing.",
	Run:   runFakedb,
}

func init() {
	rootCmd.AddCommand(fakedbCmd)
	fakedbCmd.Flags().IntVar(&fakedbJumps, "jumps", 2000, "Number of jumps to generate")
	fakedbCmd.Flags().StringVar(&fakedbOutput, "output", "skybook.db", "Path to the output SQLite database")
}

func runFakedb(cmd *cobra.Command, args []string) {
	log := slog.Default()

	if fakedbJumps <= 0 {
		log.Error("Jumps must be > 0")
		os.Exit(1)
	}

	log.Info("Starting fakedb generation", "jumps", fakedbJumps, "output", fakedbOutput)
	start := time.Now()

	// Override config completely
	config := common.NewConfig()
	config.Database.Path = fakedbOutput

	// Always start fresh — delete the existing file so we don't hit UNIQUE constraint
	// failures when the DB already contains jumps.
	if _, statErr := os.Stat(fakedbOutput); statErr == nil {
		log.Info("Removing existing database", "path", fakedbOutput)
		if err := os.Remove(fakedbOutput); err != nil {
			log.Error("Failed to remove existing database", "error", err)
			os.Exit(1)
		}
	}

	backend, err := metadata.NewBackend(config.Database, log)
	if err != nil {
		log.Error("Failed to init database", "error", err)
		os.Exit(1)
	}
	defer backend.Shutdown()

	// Seed random
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	jumps := make([]*common.Jump, 0, fakedbJumps)

	// We generate data chronologically.
	// To make it finish near present day, we calculate the required time span.
	// Avg jumps per active day = 5.5. Avg gap between active days = 10 days.
	totalSpanDays := int((float64(fakedbJumps) / 5.5) * 10.0)
	startTime := time.Now().AddDate(0, 0, -totalSpanDays)
	currentDate := common.NewDateOnly(startTime.Year(), startTime.Month(), startTime.Day())

	// Arrays for realistic generation
	homeDZ := "Skydive Local"
	destinationDZs := []string{"Skydive Algarve", "Skydive Perris", "Skydive Empuriabrava", "Skydive Dubai", "Skydive Sebastian", "Skydive Arizona"}
	aircraftList := []string{"Twin Otter", "Super Otter", "Caravan", "Skyvan", "PAC 750", "Pilatus Porter"}

	// Helper to get random item
	randString := func(items []string) string {
		return items[r.Intn(len(items))]
	}

	// Helper to get uint pointer
	uintPtr := func(v uint) *uint {
		return &v
	}

	totalGenerated := 0

	for totalGenerated < fakedbJumps {
		// Advance by 1 to 20 days to find the next active jumping day (average 10 days)
		daysToNextAct := r.Intn(20) + 1
		currentDate = currentDate.AddDays(daysToNextAct)

		// Determine if local DZ or weekend event
		isLocal := r.Float32() < 0.8
		dz := homeDZ
		if !isLocal {
			dz = randString(destinationDZs)
		}

		plane := randString(aircraftList)

		// 3 to 8 jumps per day
		jumpsToday := r.Intn(6) + 3

		for i := 0; i < jumpsToday; i++ {
			if totalGenerated >= fakedbJumps {
				break
			}

			// Discipline mix
			var jumpType common.JumpType
			var altitude, freefall uint

			prob := r.Float32()
			if prob < 0.55 {
				jumpType = common.JumpTypeFF
				altitude = dropAltitude(r)
				freefall = uint(40 + r.Intn(15))
			} else if prob < 0.70 {
				jumpType = common.JumpTypeWingsuit
				altitude = dropAltitude(r)
				freefall = uint(60 + r.Intn(30))
			} else if prob < 0.78 {
				jumpType = common.JumpTypeFS
				altitude = dropAltitude(r)
				freefall = uint(45 + r.Intn(15))
			} else if prob < 0.84 {
				jumpType = common.JumpTypeTracking
				altitude = dropAltitude(r)
				freefall = uint(50 + r.Intn(20))
			} else if prob < 0.89 {
				jumpType = common.JumpTypeAFFI
				altitude = dropAltitude(r)
				freefall = uint(55 + r.Intn(20))
			} else if prob < 0.93 {
				jumpType = common.JumpTypeCamera
				altitude = dropAltitude(r)
				freefall = uint(50 + r.Intn(20))
			} else if prob < 0.96 {
				jumpType = common.JumpTypeCRW
				altitude = dropAltitude(r)
				freefall = uint(120 + r.Intn(60))
			} else {
				jumpType = common.JumpTypeHOP
				altitude = 5000
				freefall = uint(r.Intn(6))
			}

			// Canopy progression based on total jumps
			canopySize := uint(210)
			if totalGenerated > 1000 {
				canopySize = 103
			} else if totalGenerated > 500 {
				canopySize = 135
			} else if totalGenerated > 200 {
				canopySize = 150
			} else if totalGenerated > 50 {
				canopySize = 170
			} else if totalGenerated > 20 {
				canopySize = 190
			}

			// Packjob (25%)
			description := ""
			if r.Float32() < 0.25 {
				description = "Packjob"
			}

			// Cutaway (approx 1 or 2 per 2000 jumps = 0.05% chance)
			cutaway := r.Float32() < 0.001

			jump := &common.Jump{
				UserID:       1, // AnonymousUser
				Number:       uint(totalGenerated + 1),
				Date:         currentDate,
				Dropzone:     dz,
				Aircraft:     plane,
				JumpType:     jumpType,
				Altitude:     uintPtr(altitude),
				FreefallTime: uintPtr(freefall),
				CanopySize:   uintPtr(canopySize),
				Landing:      landingResult(r),
				Pattern:      canopyPattern(totalGenerated),
				Favorite:     r.Float32() < 0.05,
				Description:  description,
				CutAway:      cutaway,
			}

			jumps = append(jumps, jump)
			totalGenerated++

			if totalGenerated%500 == 0 {
				log.Info("Generating...", "count", totalGenerated)
			}
		}
	}

	log.Info("Saving jumps to database...")

	// Batch insert for performance
	chunkSize := 500
	for i := 0; i < len(jumps); i += chunkSize {
		end := i + chunkSize
		if end > len(jumps) {
			end = len(jumps)
		}

		tx := backend.DB().Create(jumps[i:end])
		if tx.Error != nil {
			log.Error("Failed to save jumps", "error", tx.Error)
			os.Exit(1)
		}
	}

	log.Info("FakeDB complete!",
		"total_jumps", fakedbJumps,
		"user_id", 1,
		"elapsed", time.Since(start).String(),
	)

	fmt.Printf("\nSuccess! Generated %d jumps dynamically into '%s'.\n", fakedbJumps, fakedbOutput)
	fmt.Printf("Run the server against this database with:\n\n")
	fmt.Printf("  SKYBOOK_DATABASE_PATH=%s ./server/skybook serve\n\n", fakedbOutput)
}

// dropAltitude returns a realistic drop altitude:
// 90 % of the time → 13 000 ft (most common at a single-otter DZ),
// 10 % of the time → one of 10 000 / 12 000 / 15 000 ft (special loads, big-ways, etc.).
func dropAltitude(r *rand.Rand) uint {
	if r.Float32() < 0.90 {
		return 13000
	}
	alts := [3]uint{10000, 12000, 15000}
	return alts[r.Intn(3)]
}

// landingResult returns a realistic landing type:
// 90 % Stand-up, 9 % Sliding, 1 % Off-DZ.
func landingResult(r *rand.Rand) string {
	p := r.Float32()
	if p < 0.90 {
		return "Stand-up"
	}
	if p < 0.99 {
		return "Sliding"
	}
	return "Off-DZ"
}

// canopyPattern returns a canopy approach pattern based on experience level:
// ≤200 jumps → PTU (student/novice), 201–500 → 90° (intermediate), 501+ → 270° (experienced).
func canopyPattern(jumpCount int) string {
	if jumpCount <= 200 {
		return "PTU"
	}
	if jumpCount <= 500 {
		return "90°"
	}
	return "270°"
}
